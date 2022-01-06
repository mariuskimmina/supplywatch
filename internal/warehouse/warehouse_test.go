package warehouse

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetupConnections(t *testing.T) {
	t.Parallel()
	os.Setenv("SW_UDPPORT", "4444")
	os.Setenv("SW_LISTEN_IP", "0.0.0.0")
	defer os.Unsetenv("SW_UDPPORT")
	defer os.Unsetenv("SW_LISTEN_IP")

	_, err := setupUDPConn()
	require.NoError(t, err)
	_, err = setupTCPConn()
	require.NoError(t, err)

	os.Setenv("SW_UDPPORT", "4444")
	os.Setenv("SW_LISTEN_IP", "THISISNOTANIP")

	_, err = setupUDPConn()
	require.Error(t, err)
	_, err = setupTCPConn()
	require.Error(t, err)
}
