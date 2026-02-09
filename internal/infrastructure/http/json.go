//nolint:wrapcheck,errcheck
package http

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/goccy/go-json"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// CustomJSONPb is a Marshaler which marshals/unmarshals into/from JSON
// with the "google.golang.org/protobuf/encoding/protojson" marshaler.
// It supports the full functionality of protobuf unlike JSONBuiltin.
//
// The NewDecoder method returns a DecoderWrapper, so the underlying
// *json.Decoder methods can be used.
type CustomJSONPb struct {
	protojson.MarshalOptions
	protojson.UnmarshalOptions
	// Int64AsNumber controls whether int64/uint64 fields are serialized as JSON numbers.
	// Default (false): int64/uint64 are serialized as strings (protobuf default for JavaScript compatibility)
	// true: int64/uint64 are serialized as numbers (may lose precision for values > 2^53)
	Int64AsNumber bool
}

// ContentType always returns "application/json".
func (*CustomJSONPb) ContentType(_ interface{}) string {
	return "application/json"
}

// Marshal marshals "v" into JSON.
func (j *CustomJSONPb) Marshal(v interface{}) ([]byte, error) {
	msg, ok := v.(proto.Message)
	if !ok {
		return j.marshalNonProtoField(v)
	}

	var buf bytes.Buffer
	if err := j.marshalTo(&buf, v); err != nil {
		return nil, err
	}

	result := buf.Bytes()

	// Convert int64 strings to numbers if enabled
	if j.Int64AsNumber {
		converted, err := j.convertInt64StringsToNumbers(result, msg)
		if err != nil {
			return nil, err
		}
		return converted, nil
	}

	return result, nil
}

func (j *CustomJSONPb) marshalTo(w io.Writer, v interface{}) error {
	p, ok := v.(proto.Message)
	if !ok {
		buf, err := j.marshalNonProtoField(v)
		if err != nil {
			return err
		}
		_, err = w.Write(buf)
		return err
	}
	b, err := j.MarshalOptions.Marshal(p)
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	return err
}

// protoMessageType is stored to prevent constant lookup of the same type at runtime.
var protoMessageType = reflect.TypeOf((*proto.Message)(nil)).Elem()

// marshalNonProto marshals a non-message field of a protobuf message.
// This function does not correctly marshal arbitrary data structures into JSON,
// it is only capable of marshaling non-message field values of protobuf,
// i.e. primitive types, enums; pointers to primitives or enums; maps from
// integer/string types to primitives/enums/pointers to messages.
func (j *CustomJSONPb) marshalNonProtoField(v interface{}) ([]byte, error) { //nolint:gocognit //
	if v == nil {
		return []byte("null"), nil
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return []byte("null"), nil
		}
		rv = rv.Elem()
	}

	if rv.Kind() == reflect.Slice { //nolint:nestif // This is a false positive.
		if rv.IsNil() {
			if j.EmitUnpopulated {
				return []byte("[]"), nil
			}
			return []byte("null"), nil
		}

		if rv.Type().Elem().Implements(protoMessageType) {
			var buf bytes.Buffer
			if err := buf.WriteByte('['); err != nil {
				return nil, err
			}
			for i := range rv.Len() {
				if i != 0 {
					if err := buf.WriteByte(','); err != nil {
						return nil, err
					}
				}
				if err := j.marshalTo(&buf, rv.Index(i).Interface().(proto.Message)); err != nil {
					return nil, err
				}
			}
			if err := buf.WriteByte(']'); err != nil {
				return nil, err
			}

			return buf.Bytes(), nil
		}

		if rv.Type().Elem().Implements(typeProtoEnum) {
			var buf bytes.Buffer
			if err := buf.WriteByte('['); err != nil {
				return nil, err
			}
			for i := range rv.Len() {
				if i != 0 {
					if err := buf.WriteByte(','); err != nil {
						return nil, err
					}
				}
				var err error
				if j.UseEnumNumbers {
					_, err = buf.WriteString(strconv.FormatInt(rv.Index(i).Int(), 10))
				} else {
					_, err = buf.WriteString("\"" + rv.Index(i).Interface().(protoEnum).String() + "\"")
				}
				if err != nil {
					return nil, err
				}
			}
			if err := buf.WriteByte(']'); err != nil {
				return nil, err
			}

			return buf.Bytes(), nil
		}
	}

	if rv.Kind() == reflect.Map {
		m := make(map[string]*json.RawMessage)
		for _, k := range rv.MapKeys() {
			buf, err := j.Marshal(rv.MapIndex(k).Interface())
			if err != nil {
				return nil, err
			}
			m[fmt.Sprintf("%v", k.Interface())] = (*json.RawMessage)(&buf)
		}
		if j.Indent != "" {
			return json.MarshalIndent(m, "", j.Indent)
		}
		return json.Marshal(m)
	}
	if enum, ok := rv.Interface().(protoEnum); ok && !j.UseEnumNumbers {
		return json.Marshal(enum.String())
	}
	return json.Marshal(rv.Interface())
}

// Unmarshal unmarshals JSON "data" into "v".
func (j *CustomJSONPb) Unmarshal(data []byte, v interface{}) error {
	return unmarshalCustomJSONPb(data, j.UnmarshalOptions, v)
}

// NewDecoder returns a Decoder which reads JSON stream from "r".
func (j *CustomJSONPb) NewDecoder(r io.Reader) runtime.Decoder {
	d := json.NewDecoder(r)
	return DecoderWrapper{
		Decoder:          d,
		UnmarshalOptions: j.UnmarshalOptions,
	}
}

// DecoderWrapper is a wrapper around a *json.Decoder that adds
// support for protos to the Decode method.
type DecoderWrapper struct {
	*json.Decoder
	protojson.UnmarshalOptions
}

// Decode wraps the embedded decoder's Decode method to support
// protos using a jsonpb.Unmarshaler.
func (d DecoderWrapper) Decode(v interface{}) error {
	return decodeCustomJSONPb(d.Decoder, d.UnmarshalOptions, v)
}

// NewEncoder returns an Encoder which writes JSON stream into "w".
func (j *CustomJSONPb) NewEncoder(w io.Writer) runtime.Encoder {
	return runtime.EncoderFunc(func(v interface{}) error {
		if err := j.marshalTo(w, v); err != nil {
			return err
		}
		// mimic json.Encoder by adding a newline (makes output
		// easier to read when it contains multiple encoded items)
		_, err := w.Write(j.Delimiter())
		return err
	})
}

func unmarshalCustomJSONPb(data []byte, unmarshaler protojson.UnmarshalOptions, v interface{}) error {
	d := json.NewDecoder(bytes.NewReader(data))
	return decodeCustomJSONPb(d, unmarshaler, v)
}

func decodeCustomJSONPb(d *json.Decoder, unmarshaler protojson.UnmarshalOptions, v interface{}) error {
	p, ok := v.(proto.Message)
	if !ok {
		return decodeNonProtoField(d, unmarshaler, v)
	}

	// Decode into bytes for marshalling
	var b json.RawMessage
	if err := d.Decode(&b); err != nil {
		return err
	}

	return unmarshaler.Unmarshal([]byte(b), p)
}

func decodeNonProtoField(d *json.Decoder, unmarshaler protojson.UnmarshalOptions, v interface{}) error { //nolint:gocognit //
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("%T is not a pointer", v)
	}
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		if rv.Type().ConvertibleTo(typeProtoMessage) {
			// Decode into bytes for marshalling
			var b json.RawMessage
			if err := d.Decode(&b); err != nil {
				return err
			}

			return unmarshaler.Unmarshal([]byte(b), rv.Interface().(proto.Message))
		}
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Map { //nolint:nestif // This is a false positive.
		if rv.IsNil() {
			rv.Set(reflect.MakeMap(rv.Type()))
		}
		conv, ok := convFromType[rv.Type().Key().Kind()]
		if !ok {
			return fmt.Errorf("unsupported type of map field key: %v", rv.Type().Key())
		}

		m := make(map[string]*json.RawMessage)
		if err := d.Decode(&m); err != nil {
			return err
		}
		for k, v := range m {
			result := conv.Call([]reflect.Value{reflect.ValueOf(k)})
			if err := result[1].Interface(); err != nil {
				return err.(error)
			}
			bk := result[0]
			bv := reflect.New(rv.Type().Elem())
			if v == nil {
				null := json.RawMessage("null")
				v = &null
			}
			if err := unmarshalCustomJSONPb([]byte(*v), unmarshaler, bv.Interface()); err != nil {
				return err
			}
			rv.SetMapIndex(bk, bv.Elem())
		}
		return nil
	}
	if rv.Kind() == reflect.Slice { //nolint:nestif // This is a false positive.
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			var sl []byte
			if err := d.Decode(&sl); err != nil {
				return err
			}
			if sl != nil {
				rv.SetBytes(sl)
			}
			return nil
		}

		var sl []json.RawMessage
		if err := d.Decode(&sl); err != nil {
			return err
		}
		if sl != nil {
			rv.Set(reflect.MakeSlice(rv.Type(), 0, 0))
		}
		for _, item := range sl {
			bv := reflect.New(rv.Type().Elem())
			if err := unmarshalCustomJSONPb([]byte(item), unmarshaler, bv.Interface()); err != nil {
				return err
			}
			rv.Set(reflect.Append(rv, bv.Elem()))
		}
		return nil
	}
	if _, ok := rv.Interface().(protoEnum); ok {
		var repr interface{}
		if err := d.Decode(&repr); err != nil {
			return err
		}
		switch v := repr.(type) {
		case string:
			// TODO(yugui) Should use proto.StructProperties?
			return fmt.Errorf("unmarshalling of symbolic enum %q not supported: %T", repr, rv.Interface())
		case float64:
			rv.Set(reflect.ValueOf(int32(v)).Convert(rv.Type()))
			return nil
		default:
			return fmt.Errorf("cannot assign %#v into Go type %T", repr, rv.Interface())
		}
	}
	return d.Decode(v)
}

type protoEnum interface {
	fmt.Stringer
	EnumDescriptor() ([]byte, []int)
}

var typeProtoEnum = reflect.TypeOf((*protoEnum)(nil)).Elem()

var typeProtoMessage = reflect.TypeOf((*proto.Message)(nil)).Elem()

// Delimiter for newline encoded JSON streams.
func (j *CustomJSONPb) Delimiter() []byte {
	return []byte("\n")
}

var convFromType = map[reflect.Kind]reflect.Value{ //nolint:exhaustive // This is a false positive.
	reflect.String:  reflect.ValueOf(runtime.String),
	reflect.Bool:    reflect.ValueOf(runtime.Bool),
	reflect.Float64: reflect.ValueOf(runtime.Float64),
	reflect.Float32: reflect.ValueOf(runtime.Float32),
	reflect.Int64:   reflect.ValueOf(runtime.Int64),
	reflect.Int32:   reflect.ValueOf(runtime.Int32),
	reflect.Uint64:  reflect.ValueOf(runtime.Uint64),
	reflect.Uint32:  reflect.ValueOf(runtime.Uint32),
	reflect.Slice:   reflect.ValueOf(runtime.Bytes),
}

// getFieldName returns the field name to use in JSON based on UseProtoNames setting.
func getFieldName(fd protoreflect.FieldDescriptor, useProtoNames bool) string {
	if useProtoNames {
		return string(fd.Name())
	}
	return string(fd.JSONName())
}

// isInt64Kind returns true if the field kind is an int64-like type.
func isInt64Kind(kind protoreflect.Kind) bool {
	switch kind {
	case protoreflect.Int64Kind, protoreflect.Uint64Kind, protoreflect.Sint64Kind,
		protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:
		return true
	default:
		return false
	}
}

// convertFieldsRecursive recursively converts int64 string fields to numbers in the data structure.
func convertFieldsRecursive(data interface{}, md protoreflect.MessageDescriptor, useProtoNames bool) interface{} { //nolint:gocognit,gocyclo // Complex type conversion logic
	switch v := data.(type) {
	case map[string]interface{}:
		// Process map fields
		result := make(map[string]interface{})
		fields := md.Fields()
		for key, value := range v {
			// Find the field descriptor for this key
			var fd protoreflect.FieldDescriptor
			for i := 0; i < fields.Len(); i++ {
				f := fields.Get(i)
				if getFieldName(f, useProtoNames) == key {
					fd = f
					break
				}
			}

			if fd == nil {
				// Field not found in descriptor, keep as-is
				result[key] = value
				continue
			}

			// Handle repeated fields
			if fd.IsList() {
				if slice, ok := value.([]interface{}); ok {
					convertedSlice := make([]interface{}, len(slice))
					for i, item := range slice {
						if fd.Kind() == protoreflect.MessageKind {
							// Nested message
							convertedSlice[i] = convertFieldsRecursive(item, fd.Message(), useProtoNames)
						} else if isInt64Kind(fd.Kind()) {
							// int64 field in array
							convertedSlice[i] = convertStringToNumber(item)
						} else {
							convertedSlice[i] = item
						}
					}
					result[key] = convertedSlice
				} else {
					result[key] = value
				}
				continue
			}

			// Handle single fields
			if fd.Kind() == protoreflect.MessageKind {
				// Nested message
				result[key] = convertFieldsRecursive(value, fd.Message(), useProtoNames)
			} else if isInt64Kind(fd.Kind()) {
				// int64 field - convert string to number
				result[key] = convertStringToNumber(value)
			} else {
				result[key] = value
			}
		}
		return result

	case []interface{}:
		// Process array (shouldn't happen at top level for proto messages, but handle it anyway)
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = convertFieldsRecursive(item, md, useProtoNames)
		}
		return result

	default:
		return data
	}
}

// convertStringToNumber converts a string value to a number if possible.
func convertStringToNumber(value interface{}) interface{} {
	str, ok := value.(string)
	if !ok {
		return value
	}

	// Try to parse as int64 first
	if num, err := strconv.ParseInt(str, 10, 64); err == nil {
		return num
	}

	// Try to parse as uint64
	if num, err := strconv.ParseUint(str, 10, 64); err == nil {
		return num
	}

	// If parsing fails, return as-is
	return value
}

// convertInt64StringsToNumbers uses proto reflection to convert int64/uint64 fields
// from strings to numbers in the JSON output.
func (j *CustomJSONPb) convertInt64StringsToNumbers(jsonBytes []byte, msg proto.Message) ([]byte, error) {
	// Parse JSON into map
	var data map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON for int64 conversion: %w", err)
	}

	// Get message descriptor
	md := msg.ProtoReflect().Descriptor()

	// Convert fields recursively
	converted := convertFieldsRecursive(data, md, j.UseProtoNames)

	// Marshal back to JSON
	result, err := json.Marshal(converted)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal converted JSON: %w", err)
	}

	return result, nil
}
