package cycleanalyst

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ihowson/eMotoDashboard/model"
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
	}
	return f
}

func (ca *CycleAnalyst3Mock) Run(ctx context.Context) error {
	data, err := os.ReadFile(ca.File)
	if err != nil {
		return fmt.Errorf("ReadFile: %w", err)
	}

	lines := strings.Split(string(data), "\n")

	publishInterval := 100 * time.Millisecond
	nextPublish := time.Now()

	for ctx.Err() != nil { // repeatedly publish the same data in a loop
		for _, row := range lines {
			// TODO: exit loop when cancel context is Done

			time.Sleep(time.Until(nextPublish))

			dr := parseLine(row)
			publish(dr, ca.Model)

			nextPublish = nextPublish.Add(publishInterval)
		}
	}

	return nil
}
