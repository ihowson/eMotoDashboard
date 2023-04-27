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
		var value uint16
		value, err := ss.modbus.ReadRegister(address, modbus.HOLDING_REGISTER)
		if err != nil {
			continue
		}

		log.Printf("address %d value: %d", address, value)
	}
}

type datalog struct {
	Time         time.Time
	SystemStatus uint16
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

	// Dump and update the controller configuration.
	log.Printf("SABVOTON CONFIG")
	for _, conf := range DesiredConfig {
		val, err := ss.ReadUInt16(conf.Register)
		if err != nil {
			log.Printf("could not read initial %s: %v", conf.Name, err)
			continue
		}

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
		systemStatus, err := ss.ReadUInt16(RegisterSystemStatus)
		if err != nil {
			return fmt.Errorf("couldn't read SystemStatus")
		}
		ss.Model.Debugs.Store("SabvotonSystemStatus", systemStatus)

		if systemStatus == 0 {
			// Controller reads are failing. Quit to reconnect.
			return fmt.Errorf("SystemStatus is 0")

		}

		// Collect and store datalogs
		dl := datalog{
			Time: time.Now(),
			SystemStatus: systemStatus,
		}

		dlJSON, err := json.Marshal(dl)
		if err != nil {
			log.Printf("failed to marshal datalog: %v", err)
			continue
		}

		dataLogFile.Write(dlJSON)
		dataLogFile.WriteString("\n")
	}

	return ctx.Err()
}

func (ss *SabvotonSerial) ReadFloat(reg RegisterFloat16) (float64, error) {
	raw, err := ss.modbus.ReadRegister(reg.Address, modbus.HOLDING_REGISTER)
	if err != nil {
		return 0.0, fmt.Errorf("ReadRegister(Address=%d): %v", reg.Address, err)
	}
	return float64(raw) / float64(reg.Precision), nil
}

func (ss *SabvotonSerial) ReadUInt16(reg RegisterUInt16) (uint16, error) {
	return ss.modbus.ReadRegister(reg.Address, modbus.HOLDING_REGISTER)
}

func (ss *SabvotonSerial) ReadSInt16(reg RegisterSInt16) (int16, error) {
	val, err := ss.modbus.ReadRegister(reg.Address, modbus.HOLDING_REGISTER)
	if err != nil {
		return 0, fmt.Errorf("ReadSInt16 ReadRegister(Address=%d): %v", reg.Address, err)
	}
	return int16(val), nil
}

func (ss *SabvotonSerial) ReadUInt16Array(reg RegisterUInt16, length uint) ([]uint16, error) {
	return ss.modbus.ReadRegisters(reg.Address, uint16(length), modbus.HOLDING_REGISTER)
}
