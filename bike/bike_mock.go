//go:build mock

package bike

import (
	"github.com/ihowson/eMotoDashboard/cycleanalyst"
	"github.com/ihowson/eMotoDashboard/model"
)

func Build() (*model.Model, *Bike) {
	m := model.Model{}

	ca := &cycleanalyst.CycleAnalyst3Mock{
		File:  "cycleanalyst-sample-log.txt",
		Model: &m,
	}
	go ca.Run()

	return &m, nil // FIXME: bike not implemented
}
