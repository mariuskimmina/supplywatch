package warehouse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHTTPResponse(t *testing.T) {
	t.Parallel()
	res, err := NewHTTPResponse()
	require.NoError(t, err)
	require.Equal(t, 200, res.statuscode)
	require.Equal(t, "HTTP/1.1", res.HTTPVersion)
	require.Equal(t, "OK\r\n", res.reason)
	require.Equal(t, []byte{}, res.body)
	require.Equal(t, "\r\n", res.headers)
}

func TestSetHeader(t *testing.T) {
	t.Parallel()
	res, err := NewHTTPResponse()
	require.NoError(t, err)
	res.SetHeader("Content-Type", "application/json")
	require.Equal(t, "Content-Type:application/json\r\n\r\n", res.headers)
	res.SetHeader("Server", "NotApache")
	require.Equal(t, "Content-Type:application/json\r\nServer:NotApache\r\n\r\n", res.headers)
}

func TestResponseToBytes(t *testing.T) {
	t.Parallel()
	res, err := NewHTTPResponse()
	require.NoError(t, err)
	_, err = ResponseToBytes(res)
	require.NoError(t, err)
}
