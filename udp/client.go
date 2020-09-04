package udp

import (
	"fmt"
	"net"
)

type Client struct {
	conn *net.UDPConn
}

func NewClient(conn *net.UDPConn) *Client {
	return &Client{conn: conn}
}

func (c *Client) Send(m OutcomeMessage) error {
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", m.IP, m.Port))
	if err != nil {
		return fmt.Errorf("address resolve error: %s", err)
	}

	payload := append([]byte(m.CMD), m.Body...)

	if _, err := c.conn.WriteTo(payload, addr); err != nil {
		return fmt.Errorf("send packet error: %s", err)
	}
	return nil
}
