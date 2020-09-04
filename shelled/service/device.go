package service

import (
	"fmt"
	"time"

	"shelled-backend/shelled/db"
)

type DeviceService struct {
	ds db.DeviceStorage
	bs db.BankStorage
}

func NewDeviceService(ds db.DeviceStorage, bs db.BankStorage) *DeviceService {
	return &DeviceService{
		ds: ds,
		bs: bs,
	}
}

func (s *DeviceService) RegisterDevice(sn, ip string, port uint16) (*db.Device, error) {
	device, err := s.ds.FindBySerialNumber(sn)
	if err != nil {
		return nil, fmt.Errorf("get device by sn error: %s", err)
	}
	if device == nil {
		if device, err = s.ds.Create(sn, ip, port); err != nil {
			return nil, fmt.Errorf("create device error: %s", err)
		}
	}
	now := time.Now()
	device.LastRegisteredAt = &now

	if device, err = s.ds.Update(device); err != nil {
		return nil, fmt.Errorf("update device error: %s", err)
	}

	return device, nil
}

func (s *DeviceService) RegisterDeviceBank(
	deviceId int32, bankType db.BankType, number int, pins int) (*db.Bank, error) {

	bank, err := s.bs.FindByDeviceIdAndNumber(deviceId, number)
	if err != nil {
		return nil, fmt.Errorf("get bank by device_id and number error: %s", err)
	}
	if bank == nil {
		bank, err = s.bs.Register(deviceId, bankType, number, pins)
		if err != nil {
			return nil, fmt.Errorf("register device bank failed: %s", err)
		}
	}

	return bank, nil
}

func (s *DeviceService) PingDevice(sn string) (*db.Device, error) {
	device, err := s.ds.FindBySerialNumber(sn)
	if err != nil {
		return nil, fmt.Errorf("get device by sn error: %s", err)
	}
	if device == nil {
		return nil, fmt.Errorf("device with sn %s not found", sn)
	}

	now := time.Now()
	device.LastPingAt = &now

	if device, err = s.ds.Update(device); err != nil {
		return nil, fmt.Errorf("update device last ping error: %s", err)
	}

	return device, nil
}

func (s *DeviceService) HandleData(sn string, payload []byte) error {
	device, err := s.ds.FindBySerialNumber(sn)
	if err != nil || device == nil {
		return fmt.Errorf("device not found")
	}

	cnt := int(payload[0:1][0])
	pos := 1
	for i := 0; i < cnt; i++ {
		bankNum := int(payload[pos : pos+1][0])
		pos++
		bank, err := s.bs.FindByDeviceIdAndNumber(device.ID, bankNum)
		if err != nil || bank == nil {
			return fmt.Errorf("bank not found")
		}
		stateLen := 1
		if bank.Type == db.BankTypeADC || bank.Type == db.BankTypeDAC {
			stateLen = 2 * bank.Pins
		}
		state := make([]byte, stateLen)

		copy(state, payload[pos:pos+stateLen])

		bank.State = state

		bank, err = s.bs.Update(bank)
		if err != nil {
			return fmt.Errorf("bank update error: %s", err)
		}
		pos = pos + stateLen
	}

	return nil
}
