//go:build target

package main

import (
	"github.com/ihowson/eMotoDashboard/m/v2/cycleanalyst"
	"github.com/ihowson/eMotoDashboard/m/v2/model"
)

func BuildSystem() *model.Model {
	m := model.Model{}

	ca := &cycleanalyst.CycleAnalyst3Serial{
		DevicePath: "/dev/ttyUSB0",
		Model:      &m,
	}
	go ca.Run()

	return &m
}
