package errors

import (
	"github.com/abyssparanoia/goerr"
)

// publicMetadataValueKey is the reserved goerr.Values key under which
// response-visible metadata is stored. Kept unexported; callers must use
// the helpers below.
const publicMetadataValueKey = "_public_metadata"

// WithPublicMetadata attaches a key/value pair to the error that is exposed
// to API clients via a google.protobuf.Struct detail in the gRPC error
// response. Values may be any JSON-encodable type (string, number, bool,
// slice, map). Unlike goerr.WithValue, this namespace is explicitly public.
//
// The value is stored under a single reserved key in goerr.Values as a
// map[string]any, so multiple calls merge rather than overwrite.
func WithPublicMetadata(err *goerr.Error, key string, value any) *goerr.Error {
	meta := readPublicMetadata(err)
	if meta == nil {
		meta = make(map[string]any)
	}
	meta[key] = value
	return err.WithValue(publicMetadataValueKey, meta)
}

// PublicMetadata extracts the response-visible metadata attached to err
// (or any wrapped goerr.Error). Returns nil if none was set.
func PublicMetadata(err error) map[string]any {
	ge := goerr.Unwrap(err)
	if ge == nil {
		return nil
	}
	return readPublicMetadata(ge)
}

func readPublicMetadata(ge *goerr.Error) map[string]any {
	v, ok := ge.Values()[publicMetadataValueKey]
	if !ok {
		return nil
	}
	m, ok := v.(map[string]any)
	if !ok {
		return nil
	}
	return m
}
