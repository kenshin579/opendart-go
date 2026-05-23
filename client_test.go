package opendart

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient_RequiresAPIKey(t *testing.T) {
	_, err := NewClient("")
	require.Error(t, err)
}

func TestNewClient_WiresSubClients(t *testing.T) {
	c, err := NewClient("KEY", WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)
	assert.NotNil(t, c.Disclosure)
	assert.NotNil(t, c.Report)
}

func TestNewClientFromEnv(t *testing.T) {
	t.Setenv("OPENDART_API_KEY", "ENVKEY")
	c, err := NewClientFromEnv(WithCorpCodeCacheDir(t.TempDir()))
	require.NoError(t, err)
	assert.NotNil(t, c.Disclosure)

	t.Setenv("OPENDART_API_KEY", "")
	_, err = NewClientFromEnv()
	require.Error(t, err)
}
