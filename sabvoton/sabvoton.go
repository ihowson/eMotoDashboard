/*
Sabvoton serial (Modbus) interface.

Based on https://github.com/slothorpe/SabvotonCommandLineInterface/

There are many different variants of Sabvoton controllers. The one I'm
developing against here is a 72150 MQCON purchased from SiAECOSYS in early 2022.
Unlike most, it has a waterproof mini-DIN connector, not a mess of Molex
connectors, and it doesn't have an input for the motor temperature sensor. The
MQCON software fails to read the motor temperature. I've attempted to make this
work based on slothorpe's documentation, but I haven't tested it.
*/

package sabvoton

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/ihowson/eMotoDashboard/m/v2/model"
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

func (ss *SabvotonSerial) Run() error {
	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL:      fmt.Sprintf("rtu://%s", ss.DevicePath),
		Speed:    19200,
		DataBits: 8,
		Parity:   modbus.PARITY_ODD,
		StopBits: 1,
		Timeout:  300 * time.Millisecond,
	})
	if err != nil {
		return fmt.Errorf("modbus.NewClient: %w", err)
	}

	err = client.Open()
	if err != nil {
		return fmt.Errorf("modbus.Open: %w", err)
		// FIXME: multiple Open() attempts can be made on the same client until
		// the connection succeeds (i.e. err == nil), calling the constructor again
		// is unnecessary.
		// likewise, a client can be opened and closed as many times as needed.
	}
	defer client.Close()

	client.SetUnitId(1)

	// send init sequence
	_ = client.WriteRegister(RegisterAccelerateTime.Address, SabvotonInitCode)
	time.Sleep(50 * time.Millisecond)
	_ = client.WriteRegister(RegisterAccelerateTime.Address, SabvotonInitCode)
	time.Sleep(50 * time.Millisecond)
	err = client.WriteRegister(RegisterAccelerateTime.Address, SabvotonInitCode)
	if err != nil {
		return fmt.Errorf("failed to init Sabvoton: %w", err)
	}
	time.Sleep(50 * time.Millisecond)

	ss.modbus = client
	m := ss.Model

	// TODO: if we didn't get nil for the last error, keep trying to reconnect (and flag it as an error on the model)

	for {
		// TODO: if context is done, exit

		controllerTemperature := ss.ReadFloat(RegisterControllerTemperature, math.NaN())
		model.LockNStore(m, &m.ControllerTemperature, controllerTemperature)

		systemStatus := ss.ReadUInt16(RegisterSystemStatus, 0)
		ss.Model.Debugs.Store("SabvotonSystemStatus", systemStatus)

		motorSpeed := ss.ReadUInt16(RegisterMotorSpeed, 0xffff)
		ss.Model.Debugs.Store("SabvotonMotorSpeed", motorSpeed)

		mosfetStatus := ss.ReadUInt16(RegisterMOSFETStatus, 0xffff)
		ss.Model.Debugs.Store("SabvotonMOSFETStatus", mosfetStatus)

		batteryVoltage := ss.ReadFloat(RegisterBatteryVoltage, math.NaN())
		ss.Model.Debugs.Store("SabvotonBatteryVoltage", batteryVoltage)
	}

	return nil
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