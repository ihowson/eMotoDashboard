//go:build target

package cycleanalyst

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"go.bug.st/serial"

	"github.com/ihowson/eMotoDashboard/model"
)

type CycleAnalyst3Serial struct {
	DevicePath string
	Model      *model.Model

	lastUpdate time.Time
	lock       sync.Mutex
}

func (ca *CycleAnalyst3Serial) LastUpdate() time.Time {
	ca.lock.Lock()
	defer ca.lock.Unlock()

	return ca.lastUpdate
}

func (ca *CycleAnalyst3Serial) Run(ctx context.Context) error {
	dataLogFile, err := os.OpenFile("/tmp/cav3_datalog.jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open datalog file: %w", err)
	}

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

		ca.lock.Lock()
		ca.lastUpdate = time.Now()
		ca.lock.Unlock()

		dlJSON, err := json.Marshal(dr)
		if err != nil {
			log.Printf("failed to marshal datalog: %v", err)
			continue
		}

		dataLogFile.Write(dlJSON)
		dataLogFile.WriteString("\n")
	}
	err = scanner.Err()
	if err != nil {
		// FIXME: if this fails, how do we alert the rest of the system? perhaps retry forever and fire an alarm saying 'CA not reporting data'.
		// FIXME: here you might want to zero out our values in the model to show 'no data'
		return fmt.Errorf("scanner: %w", err)
	}

	return nil
}
