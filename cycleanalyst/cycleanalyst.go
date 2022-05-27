//go:build target

package cycleanalyst

import "context"
import "fmt"
import "go.bug.st/serial"
import "github.com/ihowson/eMotoDashboard/m/v2/model"

type CycleAnalyst3Serial struct {
	DevicePath string
	Model      *model.Model

	cancelContext context.Context
	port          *serial.Port
}

func (ca *CycleAnalyst3Serial) Run() (context.CancelFunc, error) {
	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 1,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(ca.DevicePath, mode)
	if err != nil {
		return nil, fmt.Errorf("Open: %w", err)
	}

	ca.port = &port

	go ca.loop()

	ctx, cancel := context.WithCancel(context.Background())
	ca.cancelContext = ctx
	return cancel, nil
}

func (ca *CycleAnalyst3Serial) loop() {
	// defer soething() // if we error, we sohuld quit?
	// watch the cancelContext

	// https://github.com/bugst/go-serial

	/*
		buff := make([]byte, 100)
		for {
			n, err := port.Read(buff)
			if err != nil {
				log.Fatal(err)
				break
			}
			if n == 0 {
				fmt.Println("\nEOF")
				break
			}
			fmt.Printf("%v", string(buff[:n]))
		}
	*/

}
