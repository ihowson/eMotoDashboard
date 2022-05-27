//go:build mock

package main

import (
	"github.com/ihowson/eMotoDashboard/m/v2/cycleanalyst"
	"github.com/ihowson/eMotoDashboard/m/v2/model"
)

func BuildSystem() *model.Model {
	m := model.Model{}

	ca := &cycleanalyst.CycleAnalyst3Mock{
		File:  "cycleanalyst-sample-log.txt",
		Model: &m,
	}
	go ca.Run()

	return &m
}
