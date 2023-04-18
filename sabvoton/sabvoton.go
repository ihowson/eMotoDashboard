/*
Sabvoton serial (Modbus) interface.

Based on https://github.com/slothorpe/SabvotonCommandLineInterface/

There are many different variants of Sabvoton controllers. The one I'm
developing against here is a 72150 MQCON purchased from SiAECOSYS in early 2022.
Unlike most, it has a waterproof mini-DIN connector, not a mess of Molex
connectors, and it doesn't have an input for the motor temperature sensor. The
MQCON software fails to read the motor temperature. I've attempted to make this
work based on slothorpe's documentation, but I haven't tested it.

Lots of good technical discussion at
https://www.elektroroller-forum.de/viewtopic.php?f=9&t=2671&start=370
*/

package sabvoton

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ihowson/eMotoDashboard/model"
	"github.com/simonvetter/modbus"
)

type SabvotonSerial struct {
	DevicePath string
	Model      *model.Model
	modbus     *modbus.ModbusClient
}

const SabvotonInitCode = uint16(13345)

// DumpAllValues reads all values from the controller and writes them to the
// log. It's intended for debugging.
func (ss *SabvotonSerial) DumpAllValues() {
	for address := uint16(0); address < 4096; address++ {
		// for address := uint16(2548); address < 4096; address++ {
		var value uint16
		value, err := ss.modbus.ReadRegister(address, modbus.HOLDING_REGISTER)
		if err != nil {
			// log.Printf("address %d failed: %v", address, err)
			continue
		}

		log.Printf("address %d value: %d", address, value)
	}
}

type datalog struct {
	Now                       time.Time
	SystemStatus              uint16
	MotorSpeed                uint16
	MotorAngle                uint16
	HallStatus                uint16
	MOSFETStatus              uint16
	ControllerTemperature2562 uint16
	ControllerTemperature2754 uint16
}

func (ss *SabvotonSerial) Run(ctx context.Context) error {
	dataLogFile, err := os.OpenFile("/tmp/svmc_datalog.jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open datalog file: %w", err)
	}

	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:      fmt.Sprintf("rtu://%s", ss.DevicePath),
		Speed:    19200,
		DataBits: 8,
		Parity:   modbus.PARITY_ODD, //nolint:nosnakecase
		StopBits: 1,
		Timeout:  300 * time.Millisecond,
	})
	if err != nil {
		return fmt.Errorf("modbus.NewClient: %w", err)
	}

	err = client.Open()
	if err != nil {
		return fmt.Errorf("modbus.Open: %w", err)
	}
	defer client.Close()

	client.SetUnitId(1)

	// send init sequence
	_ = client.WriteRegister(RegisterInitial.Address, SabvotonInitCode)
	time.Sleep(50 * time.Millisecond)
	_ = client.WriteRegister(RegisterInitial.Address, SabvotonInitCode)
	time.Sleep(50 * time.Millisecond)
	err = client.WriteRegister(RegisterInitial.Address, SabvotonInitCode)
	if err != nil {
		return fmt.Errorf("failed to init Sabvoton: %w", err)
	}
	time.Sleep(50 * time.Millisecond)

	ss.modbus = client

	// TODO: if we didn't get nil for the last error, keep trying to reconnect (and flag it as an error on the model)

	// Dump and update the controller configuration.
	log.Printf("SABVOTON CONFIG")
	for _, conf := range DesiredConfig {
		val := ss.ReadUInt16(conf.Register, 0)
		log.Printf("%s: %d", conf.Name, val)

		if conf.Value != val {
			log.Printf("UPDATE %s to %d", conf.Name, conf.Value)
			err := client.WriteRegister(conf.Register.Address, conf.Value)
			if err != nil {
				log.Printf("failed to update %s: %v", conf.Name, err)
			}
		}
	}

	for ctx.Err() == nil {
		systemStatus := ss.ReadUInt16(RegisterSystemStatus, 0)
		ss.Model.Debugs.Store("SabvotonSystemStatus", systemStatus)

		if systemStatus == 0 {
			// Controller reads are failing. Quit to reconnect.
			return fmt.Errorf("Sabvoton is not responding")
		}

		motorSpeed := ss.ReadUInt16(RegisterMotorSpeed, 0)
		ss.Model.Debugs.Store("RegisterMotorSpeed", motorSpeed)

		// Collect and store datalogs
		dl := datalog{
			Now:          time.Now(),
			SystemStatus: systemStatus,
			MotorSpeed:   motorSpeed,
			MotorAngle:   ss.ReadUInt16(RegisterMotorAngle, 0),
			HallStatus:   ss.ReadUInt16(RegisterHallStatus, 0),
			MOSFETStatus: ss.ReadUInt16(RegisterMOSFETStatus, 0),
		}

		dlJSON, err := json.Marshal(dl)
		if err != nil {
			log.Printf("failed to marshal datalog: %v", err)
			continue
		}

		dataLogFile.Write(dlJSON)
		dataLogFile.WriteString("\n")

		// This is within 1V of the CAv3
		// batteryVoltage := ss.ReadFloat(RegisterBatteryVoltage, math.NaN())
		// ss.Model.Debugs.Store("SabvotonBatteryVoltage", batteryVoltage)

		// log.Printf("controllerTemperature=%v systemStatus=%v motorSpeed=%v mosfetStatus=%v batteryVoltage=%v", controllerTemperature, systemStatus, motorSpeed, mosfetStatus, batteryVoltage)
	}

	return ctx.Err()
}

func (ss *SabvotonSerial) ReadFloat(reg RegisterFloat16, errValue float64) float64 {
	raw, err := ss.modbus.ReadRegister(reg.Address, modbus.HOLDING_REGISTER)
	if err != nil {
		log.Printf("ReadRegister(Address=%d): %v", reg.Address, err)
		return errValue
	}
	return float64(raw) / float64(reg.Precision)
}

func (ss *SabvotonSerial) ReadUInt16(reg RegisterUInt16, errValue uint16) uint16 {
	raw, err := ss.modbus.ReadRegister(reg.Address, modbus.HOLDING_REGISTER)
	if err != nil {
		log.Printf("ReadRegister(Address=%d): %v", reg.Address, err)
		return errValue
	}
	return raw
}
