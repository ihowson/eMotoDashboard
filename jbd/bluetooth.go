package jbd

// Massive thanks to Eric Poulsen for https://gitlab.com/bms-tools/bms-tools for
// code and documentation around register and command formats.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ihowson/eMotoDashboard/model"
	"tinygo.org/x/bluetooth"
)

var ErrTimeout = errors.New("request timed out")

// JBDBluetooth sets up a Bluetooth serial link to the JBD BMS.
type JBDBluetooth struct {
	Address string // MAC address of the BMS
	Model   *model.Model

	mutex sync.Mutex

	txCharacteristic bluetooth.DeviceCharacteristic
	rxBuf            bytes.Buffer
	rxPayloads       chan []byte

	latestBasicInfo    BasicInfo
	latestCellVoltages CellVoltages
}

func (jbd *JBDBluetooth) Run(ctx context.Context) error {
	jbd.rxPayloads = make(chan []byte)

	dataLogFile, err := os.OpenFile("/tmp/bms_datalog.jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open datalog file: %w", err)
	}

	var adapter = bluetooth.DefaultAdapter
	err = adapter.Enable()
	if err != nil {
		return fmt.Errorf("enable adapter: %w", err)
	}

	// FIXME: seems to need `bluetoothctl scan on`? then it shows up just fine

	scanCh := make(chan bluetooth.ScanResult)

	go func() {
		err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
			if result.Address.String() == jbd.Address {
				log.Printf("	scan found device: %s %s", result.Address.String(), result.LocalName())
				adapter.StopScan()
				scanCh <- result
			}
		})

		if err != nil {
			log.Printf("    scan failed: %v", err)
		}
	}()

	select {
	case <-scanCh:
		break
	case <-ctx.Done():
		return ctx.Err()
	}

	log.Printf("CONNECTING")

	mac, err := bluetooth.ParseMAC(jbd.Address)
	if err != nil {
		return fmt.Errorf("parse MAC '%s': %w", jbd.Address, err)
	}

	dev, err := adapter.Connect(
		bluetooth.Address{bluetooth.MACAddress{MAC: mac}},
		bluetooth.ConnectionParams{})
	if err != nil {
		return fmt.Errorf("failed to connect to device %v: %w", jbd.Address, err)
	}
	defer func() {
		err := dev.Disconnect()
		if err != nil {
			log.Printf("failed to disconnect: %v", err)
		}
	}()

	services, err := dev.DiscoverServices([]bluetooth.UUID{bluetooth.New16BitUUID(0xff00)})
	if err != nil {
		return fmt.Errorf("failed to discover services: %w", err)
	}

	svc := services[0]

	allChars, err := svc.DiscoverCharacteristics(nil)
	log.Printf("allChars: %v", allChars)

	rxUUID := bluetooth.New16BitUUID(0xff01)
	txUUID := bluetooth.New16BitUUID(0xff02)

	rxChars, err := svc.DiscoverCharacteristics([]bluetooth.UUID{rxUUID})
	if err != nil {
		return fmt.Errorf("failed to discover rx characteristic: %w", err)
	}

	rxChar := rxChars[0]

	txChars, err := svc.DiscoverCharacteristics([]bluetooth.UUID{txUUID})
	if err != nil {
		return fmt.Errorf("failed to discover tx characteristic: %w", err)
	}

	jbd.txCharacteristic = txChars[0]

	err = rxChar.EnableNotifications(func(buf []byte) {
		if len(buf) == 127 || (len(buf) == 1 && buf[0] == 0) {
			// FIXME: we're getting the notifications constantly and I don't
			// know why. This prevents them consuming all CPU time.
			time.Sleep(50 * time.Millisecond)
			return
		}

		jbd.handleRx(buf)
	})
	if err != nil {
		return fmt.Errorf("rxChar enable notifications: %w", err)
	}

	for ctx.Err() == nil {
		info, err := jbd.ReadBasicInfo(ctx)
		if errors.Is(err, ErrTimeout) {
			continue // try again but don't reconnect until we get a hard error
		} else if err != nil {
			return fmt.Errorf("read basic info: %w", err)
		}

		jbd.mutex.Lock()
		jbd.latestBasicInfo = info
		jbd.mutex.Unlock()

		// Write to datalog.
		dlJSON, err := json.Marshal(info)
		if err != nil {
			log.Printf("failed to marshal datalog: %v", err)
			continue
		}

		dataLogFile.Write(dlJSON)
		dataLogFile.WriteString("\n")

		// TODO: if we're Running poll as fast as possible. In Charging, poll slowly.
		time.Sleep(1 * time.Second)

		cellVoltages, err := jbd.ReadCellVoltages(ctx)
		if errors.Is(err, ErrTimeout) {
			continue // try again but don't reconnect until we get a hard error
		} else if err != nil {
			return fmt.Errorf("read cell voltages: %w", err)
		}

		jbd.mutex.Lock()
		jbd.latestCellVoltages = cellVoltages
		jbd.mutex.Unlock()

		time.Sleep(1 * time.Second)
	}

	return ctx.Err()
}

func checksum(arr []byte) uint16 {
	sum := uint16(0)

	for _, b := range arr {
		sum += uint16(b)
	}

	return uint16(0x10000 - uint32(sum))
}

func (jbd *JBDBluetooth) handleRx(buf []byte) {
	// FIXME: we probably need this
	// jbd.mutex.Lock()
	// defer jbd.mutex.Unlock()

	jbd.rxBuf.Write(buf)

	jbd.checkRxBuf()
}

func (jbd *JBDBluetooth) checkRxBuf() {
	// Response frame format: [start address status length <0..n data> cs1 cs2 stop]

	// Do we have enough data for a frame?
	if jbd.rxBuf.Len() < 7 {
		return
	}

	buf := jbd.rxBuf.Bytes()

	// Is the first byte a start byte?
	if buf[0] != 0xdd {
		// frame is corrupt
		// TODO: we might abort whatever transaction is in progress since it will not succeed
		jbd.rxBuf.Reset()
		return
	}

	registerNumber := int(buf[1])
	_ = registerNumber
	success := int(buf[2])
	payloadSize := int(buf[3])

	if success != 0 {
		panic("dunno")
	}

	// Do we have all of the data?
	frameSize := payloadSize + 7
	if len(buf) < frameSize {
		// waiting for some more
		return
	}

	// Is the last byte a end byte?
	if buf[frameSize-1] != 0x77 {
		// frame is corrupt
		// TODO: we might abort whatever transaction is in progress since it will not succeed
		jbd.rxBuf.Reset()
		return
	}

	jbd.rxBuf.Next(frameSize)

	buf = buf[:frameSize]

	// TODO: Is the checksum correct?
	// if checksum(jbd.rxBuf.Bytes()[1:jbd.rxBuf.Len()-3]) != uint16(jbd.rxBuf.Bytes()[jbd.rxBuf.Len()-3])<<8|uint16(jbd.rxBuf.Bytes()[jbd.rxBuf.Len()-2]) {
	// 	jbd.rxBuf.Reset()
	// 	return
	// }

	// We have a valid frame
	payload := buf[4 : frameSize-3]
	jbd.rxBuf.Reset()
	jbd.rxPayloads <- payload
}

func packetize(payload []byte) []byte {
	cs := checksum(payload[1:]) // ignore the command byte

	buf := bytes.Buffer{}
	buf.WriteByte(0xdd) // start byte
	buf.Write(payload)
	buf.WriteByte(byte((cs & 0xff00) >> 8))
	buf.WriteByte(byte((cs & 0xff)))
	buf.WriteByte(0x77) // end byte

	return buf.Bytes()
}

func writeRequest(address byte, data []byte) []byte {
	return commandRequest(address, 0x5a, data)
}

func readRequest(address byte, data []byte) []byte {
	return commandRequest(address, 0xa5, data)
}

func commandRequest(address byte, command byte, data []byte) []byte {
	buf := bytes.Buffer{}
	buf.WriteByte(command)
	buf.WriteByte(address)
	buf.WriteByte(byte(len(data)))
	buf.Write(data)

	return packetize(buf.Bytes())
}

func (jbd *JBDBluetooth) RawRequest(ctx context.Context, frame []byte) ([]byte, error) {
	// TODO: do we need this?
	// jbd.mutex.Lock()
	// defer jbd.mutex.Unlock()

	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	n, err := jbd.txCharacteristic.WriteWithoutResponse(frame)
	_ = n
	if err != nil {
		return nil, fmt.Errorf("WriteWithoutResponse: %w", err)
	}

	// Wait until we get a response.
	select {
	case response := <-jbd.rxPayloads:
		return response, nil
	case <-reqCtx.Done():
		return nil, ErrTimeout
	}
}
