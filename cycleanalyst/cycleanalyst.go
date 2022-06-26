//go:build target

package cycleanalyst

import (
	"bufio"
	"context"
	"fmt"

	"go.bug.st/serial"

	"github.com/ihowson/eMotoDashboard/m/v2/model"
)

type CycleAnalyst3Serial struct {
	DevicePath string
	Model      *model.Model
}

func (ca *CycleAnalyst3Serial) Run(ctx context.Context) error {
	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(ca.DevicePath, mode)
	if err != nil {
		return fmt.Errorf("Open(%s): %w", ca.DevicePath, err)
	}
	defer port.Close()
	// TODO: watch the cancelContext

	scanner := bufio.NewScanner(port)
	for scanner.Scan() {
		dr := parseLine(scanner.Text())
		publish(dr, ca.Model)
	}
	err = scanner.Err()
	if err != nil {
		// FIXME: if this fails, how do we alert the rest of the system? perhaps retry forever and fire an alarm saying 'CA not reporting data'.
		// FIXME: here you might want to zero out our values in the model to show 'no data'
		return fmt.Errorf("scanner: %w", err)
	}

	return nil
}
