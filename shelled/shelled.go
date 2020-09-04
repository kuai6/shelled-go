package shelled

import (
	"context"
	"log"
	"sort"

	"github.com/google/uuid"

	"shelled-backend/shelled/db"
	"shelled-backend/shelled/service"
)

type Event interface {
	Name() EventName
}

type BasicEvent struct {
	name EventName
	ip   string
	port uint16
	sn   string
}

func (be *BasicEvent) Name() EventName {
	return be.name
}

func (be *BasicEvent) IP() string {
	return be.ip
}

func (be *BasicEvent) Port() uint16 {
	return be.port
}

func (be *BasicEvent) SN() string {
	return be.sn
}

type Result interface {
	Error() error
	IsSuccess() bool
	Result() []byte
}

type DispatchResult struct {
}

func (dr *DispatchResult) Error() error {
	return nil
}

func (dr *DispatchResult) IsSuccess() bool {
	return true
}

func (dr *DispatchResult) Result() []byte {
	return nil
}

type Application interface {
	Init(db.DeviceStorage, db.BankStorage) error
	Dispatch(Event)
	Destroy() error
	RegisterHandler(EventName, Handler) Application
}

type Handler interface {
	Handle(context.Context, Event) error
	Order() int
}

type Shelled struct {
	handlers map[EventName][]Handler
}

func (a *Shelled) Init(
	ds db.DeviceStorage,
	bs db.BankStorage) error {
	a.handlers = make(map[EventName][]Handler)

	deviceService := service.NewDeviceService(ds, bs)

	a.RegisterHandler(DeviceRegisterEventName, NewDeviceRegisterEventHandler(100, deviceService))
	a.RegisterHandler(DevicePingEventName, NewDevicePingEventHandler(100, deviceService))

	return nil
}

func (a *Shelled) Dispatch(e Event) {
	var handlers []Handler
	var ok bool
	if handlers, ok = a.handlers[e.Name()]; !ok {
		log.Printf("dispatch error: cant find handler for event: %s", e.Name())
		return
	}
	requestId, err := uuid.NewUUID()
	if err != nil {
		log.Printf("dispatch error: cant generate new requestId: %s", err)
		return
	}

	go func() {
		ctx := context.WithValue(context.Background(), "requestID", requestId)
		defer ctx.Done()
		for _, h := range handlers {
			if err := h.Handle(ctx, e); err != nil {
				log.Printf("dispatch event \"%s\" error: %s", e.Name(), err)
			}
		}
	}()
}

func (a *Shelled) Destroy() error {
	return nil
}

func (a *Shelled) RegisterHandler(name EventName, h Handler) Application {
	a.handlers[name] = append(a.handlers[name], h)

	sort.Slice(a.handlers[name], func(i, j int) bool {
		return a.handlers[name][i].Order() < a.handlers[name][j].Order()
	})

	return a
}
