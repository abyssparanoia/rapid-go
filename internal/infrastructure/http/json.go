//nolint:wrapcheck,errcheck
package http

import (
	"bytes"
	"fmt"
	"io"
	"math"
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
// It is a copy of grpc-gateway's runtime.JSONPb that uses goccy/go-json for
// non-proto values and adds the Int64AsNumber option. Keep the shared parts in
// sync with runtime/marshal_jsonpb.go when upgrading grpc-gateway.
//
// The NewDecoder method returns a DecoderWrapper, so the underlying
// *json.Decoder methods can be used.
type CustomJSONPb struct {
	protojson.MarshalOptions
	protojson.UnmarshalOptions
	// Int64AsNumber controls whether 64-bit integer fields (int64, uint64, sint64,
	// fixed64, sfixed64, google.protobuf.Int64Value, google.protobuf.UInt64Value)
	// are serialized as JSON numbers, including inside repeated and map fields.
	// Default (false): they are serialized as strings (protobuf default for JavaScript compatibility)
	// true: they are serialized as numbers (may lose precision for values > 2^53)
	Int64AsNumber bool
}

// ContentType always returns "application/json".
func (*CustomJSONPb) ContentType(_ interface{}) string {
	return "application/json"
}

// Marshal marshals "v" into JSON.
func (j *CustomJSONPb) Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := j.marshalTo(&buf, v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (j *CustomJSONPb) marshalTo(w io.Writer, v interface{}) error {
	p, ok := v.(proto.Message)
	if !ok {
		buf, err := j.marshalNonProtoField(v)
		if err != nil {
			return err
		}
		if j.Indent != "" {
			b := &bytes.Buffer{}
			if err = json.Indent(b, buf, "", j.Indent); err != nil {
				return err
			}
			buf = b.Bytes()
		}
		_, err = w.Write(buf)
		return err
	}
	b, err := j.MarshalOptions.Marshal(p)
	if err != nil {
		return err
	}
	if j.Int64AsNumber {
		b, err = j.convertInt64StringsToNumbers(b, p)
		if err != nil {
			return err
		}
	}

	_, err = w.Write(b)
	return err
}

// protoMessageType is stored to prevent constant lookup of the same type at runtime.
var protoMessageType = reflect.TypeFor[proto.Message]()

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
	for rv.Kind() == reflect.Pointer {
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

func decodeNonProtoField(d *json.Decoder, unmarshaler protojson.UnmarshalOptions, v interface{}) error { //nolint:gocognit,gocyclo,cyclop,funlen // mirrors grpc-gateway runtime.decodeNonProtoField
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		return fmt.Errorf("%T is not a pointer", v)
	}
	for rv.Kind() == reflect.Pointer {
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

		var raw json.RawMessage
		if err := d.Decode(&raw); err != nil {
			return err
		}
		m := make(map[string]*json.RawMessage)
		if err := json.Unmarshal(raw, &m); err != nil {
			return err
		}
		// json keeps the last value of a repeated key, while protojson
		// rejects a duplicate map key, also when two spellings convert to the
		// same key (e.g. "1" and "01"). Reject both here too.
		if k, ok := duplicateObjectKey(raw); ok {
			return fmt.Errorf("duplicate map key %q", k)
		}
		seen := make(map[any]struct{}, len(m))
		for k, v := range m {
			result := conv.Call([]reflect.Value{reflect.ValueOf(k)})
			if err := result[1].Interface(); err != nil {
				return err.(error)
			}
			bk := result[0]
			if _, ok := seen[bk.Interface()]; ok {
				return fmt.Errorf("duplicate map key %q", k)
			}
			seen[bk.Interface()] = struct{}{}
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
			return fmt.Errorf("unmarshaling of symbolic enum %q not supported: %T", repr, rv.Interface())
		case float64:
			// Enum values are int32 on the wire. A fractional or out-of-range
			// number would otherwise be silently truncated/wrapped into a bogus
			// enum here, unlike protojson which rejects it, so reject it too.
			if v != math.Trunc(v) || v < math.MinInt32 || v > math.MaxInt32 {
				return fmt.Errorf("%v is not a valid value for enum %T", repr, rv.Interface())
			}
			rv.Set(reflect.ValueOf(int32(v)).Convert(rv.Type()))
			return nil
		default:
			return fmt.Errorf("cannot assign %#v into Go type %T", repr, rv.Interface())
		}
	}
	return d.Decode(v)
}

// duplicateObjectKey returns the first key that appears more than once in the
// JSON object in data.
func duplicateObjectKey(data []byte) (string, bool) {
	d := json.NewDecoder(bytes.NewReader(data))
	if t, err := d.Token(); err != nil || t != json.Delim('{') {
		return "", false
	}
	seen := make(map[string]struct{})
	for d.More() {
		t, err := d.Token()
		if err != nil {
			return "", false
		}
		k, _ := t.(string)
		if _, ok := seen[k]; ok {
			return k, true
		}
		seen[k] = struct{}{}
		var v json.RawMessage
		if err := d.Decode(&v); err != nil {
			return "", false
		}
	}
	return "", false
}

type protoEnum interface {
	fmt.Stringer
	EnumDescriptor() ([]byte, []int)
}

var typeProtoEnum = reflect.TypeFor[protoEnum]()

var typeProtoMessage = reflect.TypeFor[proto.Message]()

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

// convertInt64StringsToNumbers uses proto reflection to convert 64-bit integer
// fields from strings to numbers in the protojson output.
// Numbers already present in the output are preserved verbatim (json.Number),
// so no float64 round-trip can alter them.
func (j *CustomJSONPb) convertInt64StringsToNumbers(jsonBytes []byte, msg proto.Message) ([]byte, error) {
	d := json.NewDecoder(bytes.NewReader(jsonBytes))
	d.UseNumber()
	var data interface{}
	if err := d.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON for int64 conversion: %w", err)
	}

	converted := convertInt64Message(data, msg.ProtoReflect().Descriptor(), j.UseProtoNames)

	var (
		result []byte
		err    error
	)
	if j.Indent != "" {
		result, err = json.MarshalIndent(converted, "", j.Indent)
	} else {
		result, err = json.Marshal(converted)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to marshal converted JSON: %w", err)
	}

	return result, nil
}

const (
	wellKnownTypesPackage protoreflect.FullName = "google.protobuf"
	int64ValueFullName    protoreflect.FullName = "google.protobuf.Int64Value"
	uint64ValueFullName   protoreflect.FullName = "google.protobuf.UInt64Value"
)

// convertInt64Message converts the 64-bit integer fields of a message encoded by protojson.
func convertInt64Message(data interface{}, md protoreflect.MessageDescriptor, useProtoNames bool) interface{} {
	if md.FullName().Parent() == wellKnownTypesPackage {
		// protojson gives well-known types a custom representation (a bare value
		// for wrappers / Timestamp / Duration, free-form objects for Struct / Any)
		// that the field descriptors do not describe, so only the 64-bit wrappers
		// are converted and the rest is left untouched.
		switch md.FullName() {
		case int64ValueFullName, uint64ValueFullName:
			return convertStringToNumber(data)
		default:
			return data
		}
	}

	object, ok := data.(map[string]interface{})
	if !ok {
		return data
	}
	fields := md.Fields()
	result := make(map[string]interface{}, len(object))
	for key, value := range object {
		var fd protoreflect.FieldDescriptor
		if useProtoNames {
			fd = fields.ByName(protoreflect.Name(key))
		} else {
			fd = fields.ByJSONName(key)
		}
		if fd == nil {
			result[key] = value
			continue
		}
		result[key] = convertInt64Field(value, fd, useProtoNames)
	}
	return result
}

// convertInt64Field converts one field value, descending into repeated and map fields.
func convertInt64Field(value interface{}, fd protoreflect.FieldDescriptor, useProtoNames bool) interface{} {
	switch {
	case fd.IsMap():
		object, ok := value.(map[string]interface{})
		if !ok {
			return value
		}
		valueFd := fd.MapValue()
		result := make(map[string]interface{}, len(object))
		for key, item := range object {
			result[key] = convertInt64Scalar(item, valueFd, useProtoNames)
		}
		return result
	case fd.IsList():
		list, ok := value.([]interface{})
		if !ok {
			return value
		}
		result := make([]interface{}, len(list))
		for i, item := range list {
			result[i] = convertInt64Scalar(item, fd, useProtoNames)
		}
		return result
	default:
		return convertInt64Scalar(value, fd, useProtoNames)
	}
}

// convertInt64Scalar converts a single (non-repeated, non-map) value of the given field.
func convertInt64Scalar(value interface{}, fd protoreflect.FieldDescriptor, useProtoNames bool) interface{} {
	switch fd.Kind() { //nolint:exhaustive // every other kind is left as protojson emitted it
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return convertInt64Message(value, fd.Message(), useProtoNames)
	case protoreflect.Int64Kind, protoreflect.Uint64Kind, protoreflect.Sint64Kind,
		protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:
		return convertStringToNumber(value)
	default:
		return value
	}
}

// convertStringToNumber turns a decimal integer string into a JSON number,
// leaving anything else as-is.
func convertStringToNumber(value interface{}) interface{} {
	str, ok := value.(string)
	if !ok {
		return value
	}
	if _, err := strconv.ParseInt(str, 10, 64); err != nil {
		if _, err := strconv.ParseUint(str, 10, 64); err != nil {
			return value
		}
	}
	return json.Number(str)
}
