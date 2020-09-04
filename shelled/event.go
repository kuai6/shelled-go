package shelled

import (
	"context"
	"fmt"

	"shelled-backend/shelled/db"
	"shelled-backend/shelled/service"
)

type EventName string

const (
	DeviceRegisterEventName EventName = "device_register"
	DevicePingEventName     EventName = "device_ping"
)

type DeviceRegisterEvent struct {
	BasicEvent
	Banks []struct {
		Number   int
		BankType db.BankType
		Pins     int
	}
}

func NewDeviceRegisterEvent(sn, ip string, port uint16) *DeviceRegisterEvent {
	return &DeviceRegisterEvent{
		BasicEvent: BasicEvent{
			name: DeviceRegisterEventName,
			ip:   ip,
			port: port,
			sn:   sn,
		},
	}
}

func (e *DeviceRegisterEvent) Name() EventName {
	return e.name
}

type DeviceRegisterEventHandler struct {
	order         int
	deviceService *service.DeviceService
}

func NewDeviceRegisterEventHandler(
	order int,
	deviceService *service.DeviceService) *DeviceRegisterEventHandler {
	return &DeviceRegisterEventHandler{
		order:         order,
		deviceService: deviceService,
	}
}

func (h *DeviceRegisterEventHandler) Order() int {
	return h.order
}

func (h *DeviceRegisterEventHandler) Handle(ctx context.Context, e Event) error {

	var ev *DeviceRegisterEvent
	var ok bool
	if ev, ok = e.(*DeviceRegisterEvent); !ok {

	}
	d, err := h.deviceService.RegisterDevice(ev.sn, ev.ip, ev.port)
	if err != nil {
		return fmt.Errorf("device register filed: %s", err)
	}

	for _, b := range ev.Banks {
		if _, err := h.deviceService.RegisterDeviceBank(d.ID, b.BankType, b.Number, b.Pins); err != nil {
			return fmt.Errorf("device bank register filed: %s", err)
		}
	}

	fmt.Printf("[%s] device registered: %d\n", ctx.Value("requestID"), d.ID)

	return nil
}

type DevicePingEvent struct {
	name         EventName
	SerialNumber string
}

func NewDevicePingEvent() *DevicePingEvent {
	return &DevicePingEvent{
		name: DevicePingEventName,
	}
}

func (e *DevicePingEvent) Name() EventName {
	return e.name
}

type DevicePingEventHandler struct {
	order         int
	deviceService *service.DeviceService
}

func NewDevicePingEventHandler(
	order int,
	deviceService *service.DeviceService) *DevicePingEventHandler {
	return &DevicePingEventHandler{
		order:         order,
		deviceService: deviceService,
	}
}

func (h *DevicePingEventHandler) Order() int {
	return h.order
}

func (h *DevicePingEventHandler) Handle(ctx context.Context, e Event) error {
	var ev *DevicePingEvent
	var ok bool
	if ev, ok = e.(*DevicePingEvent); !ok {

	}
	d, err := h.deviceService.PingDevice(ev.SerialNumber)
	if err != nil {
		return fmt.Errorf("device ping filed: %s", err)
	}

	fmt.Printf("[%s] device ping handle: %d\n", ctx.Value("requestID"), d.ID)

	return nil
}
