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
	"context"
	"fmt"
	"log"
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

func (ss *SabvotonSerial) Run(ctx context.Context) error {
	// TODO: honor ctx.Done()

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
	// m := ss.Model

	// TODO: if we didn't get nil for the last error, keep trying to reconnect (and flag it as an error on the model)

	for ctx.Err() == nil {
		time.Sleep(time.Second)
		// No idea if this is correct. Never seen it vary from 30.
		// controllerTemperature := ss.ReadFloat(RegisterControllerTemperature, math.NaN())
		// model.LockNStore(m, &m.ControllerTemperatureCelcius, controllerTemperature)

		// This is the 23/13 code that changes when you go to flux weakening. You might want to display it as text on the dash especially if it's anomalous.
		// systemStatus := ss.ReadUInt16(RegisterSystemStatus, 0)
		// ss.Model.Debugs.Store("SabvotonSystemStatus", systemStatus)

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
