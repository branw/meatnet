package meatnet

type Event interface {
	isEvent()
}

type ProbeAddedEvent struct {
	Probe *Probe
}

func (event ProbeAddedEvent) isEvent() {}

type ProbeUpdatedEvent struct {
	Probe *Probe
}

func (event ProbeUpdatedEvent) isEvent() {}

// ProbeReadingChangedEvent indicates that new temperature value has been received
// for one of the probes. If InstantReading is false, ChangedSensors will contain
// all sensors that have different values.
type ProbeReadingChangedEvent struct {
	Probe *Probe

	InstantReading bool
	ChangedSensors []Sensor
}

func (event ProbeReadingChangedEvent) isEvent() {}

type DeviceAddedEvent struct {
	Device *Device
}

func (event DeviceAddedEvent) isEvent() {}

type DeviceUpdatedEvent struct {
	Device *Device
}

func (event DeviceUpdatedEvent) isEvent() {}
