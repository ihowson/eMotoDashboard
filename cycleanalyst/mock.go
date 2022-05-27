package cycleanalyst

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ihowson/eMotoDashboard/m/v2/model"
)

// CycleAnalyst3Mock loads a log file from disk and replays it with similar
// timing characteristics to a real Cycle Analyst. It is intended for testing.
type CycleAnalyst3Mock struct {
	File  string
	Model *model.Model
}

func parseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return math.NaN()
	} else {
		return f
	}
}

func (ca *CycleAnalyst3Mock) Run() error {
	data, err := os.ReadFile(ca.File)
	if err != nil {
		return fmt.Errorf("ReadFile: %w", err)
	}

	lines := strings.Split(string(data), "\n")

	publishInterval := time.Duration(100 * time.Millisecond)
	nextPublish := time.Now()

	for { // repeatedly publish the same data in a loop
		for _, row := range lines {
			// TODO: exit loop when cancel context is Done

			time.Sleep(time.Until(nextPublish))

			// We also tolerate # as a comment marker to simplify debugging
			if strings.HasPrefix(row, "#") {
				continue
			}

			// Split the fields.
			fields := regexp.MustCompile("[ \t]+").Split(row, -1)

			// spew.Dump(fields)
			if len(fields) < 13 {
				continue
			}

			flags := fields[13]

			preset, err := strconv.Atoi(flags[0:1])
			if err != nil {
				preset = -1
			}

			fmt.Printf("Unknown fields[12] = '%s'\n", fields[12])

			dr := CycleAnalyst3DataRow{
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

			ca.publish(dr)

			nextPublish = nextPublish.Add(publishInterval)
		}
	}

	return nil
}

func (ca *CycleAnalyst3Mock) publish(dr CycleAnalyst3DataRow) {
	// write the data into the model
	ca.Model.Lock.Lock()
	defer ca.Model.Lock.Unlock()

	ca.Model.SpeedMph = dr.Speed
	// TODO: the rest of the fields
}
