package http

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHTTPResponse(t *testing.T) {
	t.Parallel()
	res, err := NewResponse()
	require.NoError(t, err)
	require.Equal(t, 200, res.statuscode)
	require.Equal(t, "HTTP/1.1", res.HTTPVersion)
	require.Equal(t, "OK\r\n", res.reason)
	require.Equal(t, []byte{}, res.body)
	require.Equal(t, "\r\n", res.headers)
}

func TestSetHeader(t *testing.T) {
	t.Parallel()
	res, err := NewResponse()
	require.NoError(t, err)
	res.SetHeader("Content-Type", "application/json")
	require.Equal(t, "Content-Type:application/json\r\n\r\n", res.headers)
	res.SetHeader("Server", "NotApache")
	require.Equal(t, "Content-Type:application/json\r\nServer:NotApache\r\n\r\n", res.headers)
}

func TestSetStatusCode(t *testing.T) {
	t.Parallel()
	res, err := NewResponse()
	require.NoError(t, err)
	res.SetStatusCode(404)
	require.Equal(t, 404, res.statuscode)
	res.SetStatusCode(200)
	require.Equal(t, 200, res.statuscode)
}

func TestSetBody(t *testing.T) {
	t.Parallel()
	res, err := NewResponse()
	require.NoError(t, err)
	res.SetBody([]byte("TestBody"))
	require.Equal(t, []byte("TestBody"), res.body)
}

func TestSetReason(t *testing.T) {
	t.Parallel()
	res, err := NewResponse()
	require.NoError(t, err)
	res.SetReason("Not Found")
	require.Equal(t, "Not Found\r\n", res.reason)
}

func TestResponseToBytes(t *testing.T) {
	t.Parallel()
	res, err := NewResponse()
	require.NoError(t, err)
	_, err = ResponseToBytes(res)
	require.NoError(t, err)
}
