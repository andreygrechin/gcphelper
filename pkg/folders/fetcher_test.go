package folders_test

import (
	"testing"

	"github.com/andreygrechin/gcphelper/pkg/folders"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchOptions_Defaults(t *testing.T) {
	opts := folders.NewFetchOptions()

	assert.Empty(t, opts.Parent)
}

func TestFetchOptions_CustomValues(t *testing.T) {
	opts := &folders.FetchOptions{
		Parent: "organizations/123456789",
	}

	assert.Equal(t, "organizations/123456789", opts.Parent)
}

// test basic folder options functionality.
func TestClient_BasicFunctionality(t *testing.T) {
	ctx := t.Context()
	opts := &folders.FetchOptions{
		Parent: "organizations/123456789",
	}

	// basic smoke test
	require.NotNil(t, ctx)
	require.NotNil(t, opts)
	assert.Equal(t, "organizations/123456789", opts.Parent)
}
