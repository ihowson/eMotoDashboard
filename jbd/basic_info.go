package jbd

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"time"
)

type BasicInfo struct {
	Time time.Time

	PackVolts              float64
	PackAmps               float64
	PackCapacityAmpHours   float64
	DesignCapacityAmpHours float64
	Cycles                 uint
	ManufactureDate        time.Time
	CellBalancingActive    []bool
	Errors                 []BasicInfoError
	SoftwareVersion        int
	StateOfChargePercent   int
	ChargeFETConducting    bool
	DischargeFETConducting bool
	NumCells               int
	InternalTemperature    float32
	PackTemperature1       float32
	PackTemperature2       float32
}

type basicInfoRaw struct {
	PackVoltage10mV      uint16
	PackAmperes10mA      int16
	BalanceCapacity10mAh uint16
	DesignCapacity10mAh  uint16
	Cycles               uint16
	ManufactureDate      uint16 // packed fields
	CellBalancingLow     uint16
	CellBalancingHigh    uint16
	Errors               uint16
	SoftwareVersion      uint8 // this overlaps the errors field?
	StateOfCharge        uint8
	FETStatus            uint8
	NumCells             uint8
	NTCCount             uint8
}

func (jbd *JBDBluetooth) LatestBasicInfo() BasicInfo {
	jbd.mutex.Lock()
	defer jbd.mutex.Unlock()

	return jbd.latestBasicInfo
}

func (jbd *JBDBluetooth) ReadBasicInfo(ctx context.Context) (BasicInfo, error) {
	info := BasicInfo{
		Time: time.Now(),
	}

	req := readRequest(RegisterBasicInfo, []byte{})
	resp, err := jbd.RawRequest(ctx, req)
	if err != nil {
		return info, fmt.Errorf("RawRequest: %w", err)
	}

	info, err = parseBasicInfoRaw(resp)
	if err != nil {
		return info, fmt.Errorf("parseBasicInfoRaw: %w", err)
	}

	return info, nil
}

func parseBasicInfoRaw(data []byte) (BasicInfo, error) {
	info := BasicInfo{}

	sz := binary.Size(&basicInfoRaw{})

	reader := bytes.NewReader(data)
	raw := basicInfoRaw{}
	err := binary.Read(reader, binary.BigEndian, &raw)
	if err != nil {
		return info, fmt.Errorf("binary.Read: %w", err)
	}

	info = BasicInfo{
		PackVolts:              float64(raw.PackVoltage10mV) * 10.0 / 1000.0,
		PackAmps:               float64(raw.PackAmperes10mA) * 10.0 / 1000.0,
		PackCapacityAmpHours:   float64(raw.BalanceCapacity10mAh) * 10.0 / 1000.0,
		DesignCapacityAmpHours: float64(raw.DesignCapacity10mAh) * 10.0 / 1000.0,
		Cycles:                 uint(raw.Cycles),
		// TODO: missing fields
		// ManufactureDate:
		// CellBalancingActive: ,
		// Errors: []BasicInfoError{},
		SoftwareVersion:        int(raw.SoftwareVersion),
		StateOfChargePercent:   int(raw.StateOfCharge),
		ChargeFETConducting:    bool(raw.FETStatus&0x01 != 0),
		DischargeFETConducting: bool(raw.FETStatus&0x02 != 0),
		NumCells:               int(raw.NumCells),
	}

	if raw.NTCCount != 3 {
		log.Printf("WARNING: NTCCount is %d, expected 3", raw.NTCCount)
	}

	for i := 0; i < int(raw.NTCCount); i++ {
		// Stored as Kelvin * 10, apparently.
		celcius := float32(binary.BigEndian.Uint16(data[sz+i*2:2+sz+i*2])-2731) / 10.0

		switch i {
		case 0:
			info.InternalTemperature = celcius
		case 1:
			info.PackTemperature1 = celcius
		case 2:
			info.PackTemperature2 = celcius
		}
	}

	return info, nil
}
