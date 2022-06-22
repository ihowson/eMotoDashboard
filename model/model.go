package model

import (
	"sync"
)

// Abstract vehicle model
// This represents the whole 'state' of the vehicle. Hardware and mock
// interfaces mutate this.
type Model struct {
	SpeedMph                float64 // TODO: mph or kph?
	BatteryAmps             float64
	BatteryAmpHoursConsumed float64
	BatteryStateOfCharge    float64 // percent from 1.0 to 0.0
	Distance                float64 // trip miles or km TODO
	Odometer                float64 // total traveled miles or km TODO
	MotorTemperatureCelcius float64
	Gear                    string // 'N', '2', '3'
	Faults                  []string

	ControllerTemperature    float64
	BatteryVoltageCA         float64
	BatteryVoltageController float64 // TODO: difference here might measure sag in your wiring?

	FluxWeakeningActive  bool
	ControllerMotorSpeed uint16 // FIXME: not sure what this is; does it give you an RPM? or is it another speed measurement?

	Debugs sync.Map // miscellaneous values that we want to expose for debugging

	Lock sync.Mutex
}

// FIXME: expose the Debugs out somewhere -- on a page

func LockNLoad[T any](m *Model, f *T) T {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	return *f
}

func LockNStore[T any](m *Model, dest *T, newValue T) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	oldValue := *dest
	*dest = newValue

	// TODO
	_ = oldValue
	// if oldValue != newValue {
	// push to chan to notify of change
// }
}

// There are no accessor functions. In exchange for that convenience, all reads
// to fields must have Lock set or use atomic operations.
