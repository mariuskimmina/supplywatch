package warehouse

// SensorMessageCounter keep track of the number of packages we receive from each Sensor
// this way we can track how many udp packets arrive and compare it to the number
// of send udp packets on the Sensor
type SensorMessageCounter struct {
	SensorID string
	Counter  int
}

// NewSensorMessageCounter creates a new SensorMessageCounter
func NewSensorMessageCounter(id string) *SensorMessageCounter {
	return &SensorMessageCounter{
		SensorID: id,
		Counter:  1,
	}
}

// increment increments the counter for a Sensor
func (smc *SensorMessageCounter) increment() {
	smc.Counter += 1
}
