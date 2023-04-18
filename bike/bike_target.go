//go:build target

package bike

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ihowson/eMotoDashboard/cycleanalyst"
	"github.com/ihowson/eMotoDashboard/jbd"
	"github.com/ihowson/eMotoDashboard/model"
	"github.com/ihowson/eMotoDashboard/sabvoton"
	"github.com/simonvetter/modbus"
)

type Bike struct {
	CycleAnalyst *cycleanalyst.CycleAnalyst3Serial
	Sabvoton     *sabvoton.SabvotonSerial
	BMS          *jbd.JBDBluetooth
	// BDC *bdc.BDC
}

func Build() (*model.Model, *Bike) {
	m := &model.Model{}

	ctx, cancel := context.WithCancel(context.Background())
	_ = cancel // TODO: use this

	ca := &cycleanalyst.CycleAnalyst3Serial{
		DevicePath: "/dev/ttySC0",
		// WARNING: you can't use /dev/serial0 here as the core clock rate
		// varies, and the serial bitrate is a divisor off this. Effectively,
		// the bitrate varies with system load.
		Model: m,
	}
	go func() {
		for {
			log.Printf("CycleAnalyst3Serial Run")
			err := ca.Run(ctx)
			log.Printf("CycleAnalyst3Serial exited: %v", err)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	sabvoton := &sabvoton.SabvotonSerial{
		DevicePath: "/dev/ttySC1",
		Model:      m,
	}
	go func() {
		for {
			err := sabvoton.Run(ctx)
			if !errors.Is(err, modbus.ErrRequestTimedOut) {
				log.Printf("SabvotonSerial exited: %v", err)
			}
			time.Sleep(250 * time.Millisecond)
		}
	}()

	bms := &jbd.JBDBluetooth{
		// Address: "70:3e:97:08:05:4b",
		Address: "70:3E:97:08:05:4B",
		Model:   m,
	}
	go func() {
		for {
			log.Printf("JBDBluetooth Run")
			err := bms.Run(ctx)
			log.Printf("JBDBluetooth exited: %v", err)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	bike := &Bike{
		CycleAnalyst: ca,
		Sabvoton:     sabvoton,
		BMS:          bms,
		// BDC: bdc,
	}

	return m, bike
}
