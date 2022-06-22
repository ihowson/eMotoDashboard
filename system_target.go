//go:build target

package main

import (
	"github.com/ihowson/eMotoDashboard/m/v2/cycleanalyst"
	"github.com/ihowson/eMotoDashboard/m/v2/model"
	"github.com/ihowson/eMotoDashboard/m/v2/sabvoton"
)

func BuildSystem() *model.Model {
	m := model.Model{}

	ca := &cycleanalyst.CycleAnalyst3Serial{
		DevicePath: "/dev/ttySC0",
		// WARNING: you can't use /dev/serial0 here as the core clock rate
		// varies, and the serial bitrate is a divisor off this. Effectively,
		// the bitrate varies with system load.
		Model: &m,
	}
	go ca.Run()

	sabvoton := &sabvoton.SabvotonSerial{
		DevicePath: "/dev/ttyUSB0",
		Model:      &m,
	}
	go sabvoton.Run()

	return &m
}
