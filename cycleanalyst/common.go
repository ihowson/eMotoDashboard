package cycleanalyst

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ihowson/eMotoDashboard/model"
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
	// Unknown            float64 // TODO: this is the preset 'modified value'
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

func parseLine(row string) *CycleAnalyst3DataRow {
	// We also tolerate # as a comment marker to simplify debugging
	if strings.HasPrefix(row, "#") {
		return nil
	}

	// Split the fields.
	fields := regexp.MustCompile("[ \t]+").Split(row, -1)

	if len(fields) < 14 {
		log.Printf("short row '%v'", fields)
		return nil
	}

	flags := fields[13]
	preset, err := strconv.Atoi(flags[0:1])
	if err != nil {
		preset = -1
	}

	// this is the 'controlled value' for the preset -- 0 for 0% throttle, 99.9 for 99.9% throttle
	// FIXME: don't use this; use the Flags field instead
	switch fields[12] {
	case "0.00":
		preset = 1
	case "99.9":
		preset = 2
	default:
		preset = -1
	}

	return &CycleAnalyst3DataRow{
		Timestamp: time.Now(),

		AmpHours:           parseFloat(fields[0]),
		Voltage:            parseFloat(fields[1]),
		Amperes:            parseFloat(fields[2]),
		Speed:              parseFloat(fields[3]),
		Distance:           parseFloat(fields[4]),
		TemperatureCelcius: parseFloat(fields[5]),
		PASRPM:             parseFloat(fields[6]),
		HumanWatts:         parseFloat(fields[7]),
		PASTorque:          parseFloat(fields[8]),
		ThrottleInVoltage:  parseFloat(fields[9]),
		ThrottleOutVoltage: parseFloat(fields[10]),
		Acceleration:       parseFloat(fields[11]),

		// limit flags
		// 1/2/3 = Preset #
		// X = Throttle Fault
		// B = Brake
		// A = Amp Limiting
		// W = Watt Limiting
		// T = Temp Limiting
		// V = Low Volts Limiting
		// S = Speed Limiting
		// s = Low Speed Limiting
		Preset:              preset,
		ThrottleFault:       strings.Contains(flags, "X"),
		Brake:               strings.Contains(flags, "B"),
		AmpLimiting:         strings.Contains(flags, "A"),
		WattLimiting:        strings.Contains(flags, "W"),
		TemperatureLimiting: strings.Contains(flags, "T"),
		LowVoltsLimiting:    strings.Contains(flags, "V"),
		SpeedLimiting:       strings.Contains(flags, "S"),
		LowSpeedLimiting:    strings.Contains(flags, "s"),
	}
}

func publish(dr *CycleAnalyst3DataRow, model *model.Model) {
	if dr == nil || model == nil {
		log.Printf("not publishing: dr=%p model=%p\n", dr, model)
		return
	}
	// write the data into the model
	model.Lock.Lock()
	defer model.Lock.Unlock()

	model.BatteryVoltageCA = dr.Voltage
	model.SpeedMph = dr.Speed
	model.BatteryAmps = dr.Amperes
	model.BatteryAmpHoursConsumed = dr.AmpHours
	model.Distance = dr.Distance
	// model.Odometer = dr.Odometer
	model.MotorTemperatureCelcius = dr.TemperatureCelcius

	// log.Printf("preset = %v", dr.Preset)

	// TODO: the faults

	switch dr.Preset {
	case 1:
		model.Gear = "N"
	case 2:
		model.Gear = "1"
	case 3:
		model.Gear = "2"
	default:
		model.Gear = "?"
	}

	// TODO: the rest of the fields

	// TODO: you need to blank out the fields if there is no data -- find a 'mmissing daat' lib
}
