package errors

import (
	"github.com/abyssparanoia/goerr"
)

// publicMetadataValueKey is the reserved goerr.Values key under which
// response-visible metadata is stored. Kept unexported; callers must use
// the helpers below.
const publicMetadataValueKey = "_public_metadata"

// WithPublicMetadata attaches a key/value pair to the error that is intended
// to be exposed to API clients via ErrorInfo.Metadata in the gRPC error
// response. Unlike goerr.WithValue, this namespace is explicitly public.
//
// The value is stored under a single reserved key in goerr.Values as a
// map[string]string, so multiple calls merge rather than overwrite.
func WithPublicMetadata(err *goerr.Error, key, value string) *goerr.Error {
	meta := readPublicMetadata(err)
	if meta == nil {
		meta = make(map[string]string)
	}
	meta[key] = value
	return err.WithValue(publicMetadataValueKey, meta)
}

// PublicMetadata extracts the response-visible metadata attached to err
// (or any wrapped goerr.Error). Returns nil if none was set.
func PublicMetadata(err error) map[string]string {
	ge := goerr.Unwrap(err)
	if ge == nil {
		return nil
	}
	return readPublicMetadata(ge)
}

func readPublicMetadata(ge *goerr.Error) map[string]string {
	v, ok := ge.Values()[publicMetadataValueKey]
	if !ok {
		return nil
	}
	m, ok := v.(map[string]string)
	if !ok {
		return nil
	}
	return m
}
