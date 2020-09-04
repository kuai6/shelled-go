package db

import "time"

type Device struct {
	ID               int32      `db:"id"`
	Name             *string    `db:"name"`
	IP               string     `db:"ip"`
	Port             int        `db:"port"`
	SerialNumber     string     `db:"serial_number"`
	LastRegisteredAt *time.Time `db:"last_registered_at"`
	LastPingAt       *time.Time `db:"last_ping_at"`
}

type DeviceStorage interface {
	//FindOne(int32) (*Device, error)
	FindBySerialNumber(string) (*Device, error)
	Create(sn, ip string, port uint16) (*Device, error)
	Update(*Device) (*Device, error)
}
