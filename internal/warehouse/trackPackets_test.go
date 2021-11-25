package warehouse

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCounterIncrement(t *testing.T) {
	t.Parallel()
	id := uuid.New()
	smc := NewSensorMessageCounter(id)
	require.Equal(t, 1, smc.Counter)
	smc.increment()
	require.Equal(t, 2, smc.Counter)
	smc.increment()
	require.Equal(t, 3, smc.Counter)

	require.Equal(t, id, smc.SensorID)
}
