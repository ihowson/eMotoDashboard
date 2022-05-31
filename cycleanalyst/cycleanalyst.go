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

func (ca *CycleAnalyst3Serial) Run() context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			err := ca.loop(ctx, ca.Model)
			fmt.Printf("Run: %v\n", err)
		}
	}()

	return cancel
}

func (ca *CycleAnalyst3Serial) loop(ctx context.Context, model *model.Model) error {
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
		return fmt.Errorf("reading standard input:", err)
	}

	return nil
}
