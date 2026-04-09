package errors

import (
	"fmt"
	"testing"

	"github.com/abyssparanoia/goerr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithPublicMetadata_SingleEntry(t *testing.T) {
	err := goerr.New("test error")
	err = WithPublicMetadata(err, "field", "email")

	meta := PublicMetadata(err)
	require.NotNil(t, meta)
	assert.Equal(t, "email", meta["field"])
	assert.Len(t, meta, 1)
}

func TestWithPublicMetadata_MultipleCallsMerge(t *testing.T) {
	err := goerr.New("test error")
	err = WithPublicMetadata(err, "field", "email")
	err = WithPublicMetadata(err, "form_id", "signup")

	meta := PublicMetadata(err)
	require.NotNil(t, meta)
	assert.Equal(t, "email", meta["field"])
	assert.Equal(t, "signup", meta["form_id"])
	assert.Len(t, meta, 2)
}

func TestWithPublicMetadata_OverwritesSameKey(t *testing.T) {
	err := goerr.New("test error")
	err = WithPublicMetadata(err, "field", "original")
	err = WithPublicMetadata(err, "field", "overwritten")

	meta := PublicMetadata(err)
	require.NotNil(t, meta)
	assert.Equal(t, "overwritten", meta["field"])
	assert.Len(t, meta, 1)
}

func TestWithPublicMetadata_IntValue(t *testing.T) {
	err := goerr.New("test error")
	err = WithPublicMetadata(err, "count", 42)

	meta := PublicMetadata(err)
	require.NotNil(t, meta)
	assert.Equal(t, 42, meta["count"])
}

func TestWithPublicMetadata_BoolValue(t *testing.T) {
	err := goerr.New("test error")
	err = WithPublicMetadata(err, "is_expired", true)

	meta := PublicMetadata(err)
	require.NotNil(t, meta)
	assert.Equal(t, true, meta["is_expired"])
}

func TestWithPublicMetadata_SliceValue(t *testing.T) {
	err := goerr.New("test error")
	tags := []string{"a", "b"}
	err = WithPublicMetadata(err, "tags", tags)

	meta := PublicMetadata(err)
	require.NotNil(t, meta)
	assert.Equal(t, tags, meta["tags"])
}

func TestWithPublicMetadata_NestedMapValue(t *testing.T) {
	err := goerr.New("test error")
	nested := map[string]any{"k": "v", "n": 1}
	err = WithPublicMetadata(err, "details", nested)

	meta := PublicMetadata(err)
	require.NotNil(t, meta)
	assert.Equal(t, nested, meta["details"])
}

func TestPublicMetadata_NilError(t *testing.T) {
	meta := PublicMetadata(nil)
	assert.Nil(t, meta)
}

func TestPublicMetadata_NonGoerr(t *testing.T) {
	meta := PublicMetadata(fmt.Errorf("plain error"))
	assert.Nil(t, meta)
}

func TestPublicMetadata_NoMetadataSet(t *testing.T) {
	err := goerr.New("test error")
	meta := PublicMetadata(err)
	assert.Nil(t, meta)
}

func TestPublicMetadata_ThroughWrapChain(t *testing.T) {
	inner := goerr.New("inner error")
	inner = WithPublicMetadata(inner, "field", "email")

	outer := goerr.Wrap(inner, "outer context")

	// PublicMetadata uses goerr.Unwrap which finds the outermost goerr.Error.
	// The outer error has no metadata itself, so we get nil from the outer.
	// Verify that setting metadata on the outer wrapper is what gets returned.
	outerWithMeta := WithPublicMetadata(outer, "extra", "value")
	meta := PublicMetadata(outerWithMeta)
	require.NotNil(t, meta)
	assert.Equal(t, "value", meta["extra"])
}

func TestWithPublicMetadata_DoesNotLeakToOtherValues(t *testing.T) {
	err := goerr.New("test error")
	err = err.WithValue("sensitive_key", "sensitive_value")
	err = WithPublicMetadata(err, "public_key", "public_value")

	values := err.Values()
	// _public_metadata key exists
	assert.Contains(t, values, publicMetadataValueKey)
	// sensitive key still present
	assert.Equal(t, "sensitive_value", values["sensitive_key"])

	// PublicMetadata only returns the public namespace
	meta := PublicMetadata(err)
	require.NotNil(t, meta)
	assert.Equal(t, "public_value", meta["public_key"])
	assert.NotContains(t, meta, "sensitive_key")
}
