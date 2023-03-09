package jbd

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"
)

type CellVoltages struct {
	Time  time.Time
	Volts []float32
}

func (jbd *JBDBluetooth) LatestCellVoltages() CellVoltages {
	jbd.mutex.Lock()
	defer jbd.mutex.Unlock()

	return jbd.latestCellVoltages
}

func (jbd *JBDBluetooth) ReadCellVoltages(ctx context.Context) (CellVoltages, error) {
	cv := CellVoltages{
		Time: time.Now(),
	}

	req := readRequest(RegisterCellVoltages, []byte{})
	resp, err := jbd.RawRequest(ctx, req)
	if err != nil {
		return cv, fmt.Errorf("RawRequest: %w", err)
	}

	if len(resp)%2 != 0 {
		return cv, fmt.Errorf("invalid response length")
	}

	for i := 0; i < len(resp); i += 2 {
		cv.Volts = append(cv.Volts, float32(binary.BigEndian.Uint16(resp[i:i+2]))/1000.0)
	}

	return cv, nil
}
