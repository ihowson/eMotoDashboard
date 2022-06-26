//go:build target

package main

import (
	"context"
	"log"
	"time"

	"github.com/ihowson/eMotoDashboard/m/v2/cycleanalyst"
	"github.com/ihowson/eMotoDashboard/m/v2/model"
	"github.com/ihowson/eMotoDashboard/m/v2/sabvoton"
)

func BuildSystem() *model.Model {
	m := model.Model{}

	ctx, cancel := context.WithCancel(context.Background())
	_ = cancel

	ca := &cycleanalyst.CycleAnalyst3Serial{
		DevicePath: "/dev/ttySC0",
		// WARNING: you can't use /dev/serial0 here as the core clock rate
		// varies, and the serial bitrate is a divisor off this. Effectively,
		// the bitrate varies with system load.
		Model: &m,
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
		Model:      &m,
	}
	go func() {
		for {
			log.Printf("SabvotonSerial Run")
			err := sabvoton.Run(ctx)
			log.Printf("SabvotonSerial exited: %v", err)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	return &m
}
