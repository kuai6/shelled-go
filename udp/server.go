package udp

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
)

type Server struct {
	handlers map[CMD]Handler
	conn     *net.UDPConn
}

func NewServer(conn *net.UDPConn, handlers map[CMD]Handler) *Server {
	return &Server{
		conn:     conn,
		handlers: handlers,
	}
}

func (s *Server) Run(ctx context.Context, shutdown chan bool) error {
	for {
		select {
		case <-shutdown:
			log.Print("udp server: shutting down")
			return nil
		default:
			var buf [128]byte
			n, addr, err := s.conn.ReadFromUDP(buf[0:])
			if err != nil {
				log.Printf("udp server: error reading udp packet: %s", err)
			} else {
				if err := s.handle(buf[0:n], addr.IP.String(), uint16(addr.Port)); err != nil {
					log.Printf("udp server: error handling udp packet: %s", err)
				}
			}
		}
	}
}

func (s *Server) handle(b []byte, ip string, port uint16) error {
	r := bytes.Runes(b)
	if len(r) < 12 {
		return nil
	}

	m := &IncomeMessage{
		CMD:  CMD(r[0:4]),
		SN:   string(r[4:12]),
		IP:   ip,
		Port: port,
		Body: r[12:],
	}

	var h Handler
	var ok bool
	if h, ok = s.handlers[m.CMD]; ok != true {
		return fmt.Errorf("no handler found for %s commnd", m.CMD)
	}

	return h.Handle(*m)
}

func (s *Server) send(b []byte, ip string, port uint16) error {

	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return fmt.Errorf("address resolve error: %s", err)
	}

	if _, err := s.conn.WriteTo(b, addr); err != nil {
		return fmt.Errorf("send packet error: %s", err)
	}

	return nil
}
