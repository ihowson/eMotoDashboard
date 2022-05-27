package cycleanalyst

import (
	"time"
)

type CycleAnalyst3DataRow struct {
	Timestamp time.Time

	AmpHours           float64
	Voltage            float64
	Amperes            float64
	Speed              float64
	Distance           float64
	TemperatureCelcius float64
	PASRPM             float64
	HumanWatts         float64
	PASTorque          float64
	ThrottleInVoltage  float64
	ThrottleOutVoltage float64
	Acceleration       float64
	// Unknown            float64
	Preset int

	// limit flags
	ThrottleFault       bool
	Brake               bool
	AmpLimiting         bool
	WattLimiting        bool
	TemperatureLimiting bool
	LowVoltsLimiting    bool
	SpeedLimiting       bool
	LowSpeedLimiting    bool
}
