package udp

import (
	"fmt"
	"unicode/utf8"

	"shelled-backend/shelled"
	"shelled-backend/shelled/db"
)

type Handler interface {
	Handle(packet IncomeMessage) error
}

type HelloHandler struct {
	application shelled.Application
	client      *Client
}

func NewHelloHandler(application shelled.Application, client *Client) *HelloHandler {
	return &HelloHandler{
		application: application,
		client:      client,
	}
}

func (h *HelloHandler) Handle(m IncomeMessage) error {
	e := shelled.NewDeviceRegisterEvent(m.SN, m.IP, m.Port)

	// count of banks
	cnt := int(m.Body[0:1][0])

	pos := 1
	for i := 0; i < cnt; i++ {
		data := m.Body[pos : pos+6]
		if len(data) < 6 {
			return fmt.Errorf("parse hello packet error: could'nt parse bank data")
		}

		e.Banks = append(e.Banks, struct {
			Number   int
			BankType db.BankType
			Pins     int
		}{
			Number:   int(m.Body[pos+4 : pos+5][0]),
			BankType: BankToDbType(string(m.Body[pos : pos+4])),
			Pins:     int(m.Body[pos+5 : pos+6][0])})

		pos = pos + 6
	}

	h.application.Dispatch(e)

	o := OutcomeMessage{
		CMD:  CMD_CONF,
		IP:   m.IP,
		Port: m.Port,
		Body: []byte{0x88, 0x13, 0x00, 0x00, 0xf4, 0x01, 0x00, 0x00},
	}

	if err := h.client.Send(o); err != nil {
		return fmt.Errorf("hello handle erorr: can't send conf: %s", err)
	}

	return nil
}

type PingHandler struct {
	application shelled.Application
	client      *Client
}

func NewPingHandler(application shelled.Application, client *Client) *PingHandler {
	return &PingHandler{
		application: application,
		client:      client,
	}
}

func (h *PingHandler) Handle(m IncomeMessage) error {
	e := shelled.NewDevicePingEvent(m.SN, m.IP, m.Port)

	h.application.Dispatch(e)

	o := OutcomeMessage{
		CMD:  CMD_PONG,
		IP:   m.IP,
		Port: m.Port,
	}

	if err := h.client.Send(o); err != nil {
		return fmt.Errorf("ping handle erorr: can't send pong: %s", err)
	}

	return nil
}

type DataHandler struct {
	application shelled.Application
}

func NewDataHandler(application shelled.Application) *DataHandler {
	return &DataHandler{
		application: application,
	}
}

func (h *DataHandler) Handle(m IncomeMessage) error {
	e := shelled.NewDeviceDataEvent(m.SN, m.IP, m.Port)

	e.Payload = make([]byte, len(m.Body)*utf8.UTFMax)

	count := 0
	for _, r := range m.Body {
		count += utf8.EncodeRune(e.Payload[count:], r)
	}
	e.Payload = e.Payload[:count]

	h.application.Dispatch(e)

	return nil
}
