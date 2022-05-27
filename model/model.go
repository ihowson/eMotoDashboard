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

	Lock sync.Mutex
}

// func (m *Model) LockNLoadFloat64(f *float64) float64 {
// 	m.Lock.Lock()
// 	defer m.Lock.Unlock()
// 	return *f
// }

// There are no accessor functions. In exchange for that convenience, all reads
// to fields must have Lock set or use atomic operations.
