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
	return nil, nil
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
