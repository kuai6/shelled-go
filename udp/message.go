package udp

import (
	"shelled-backend/shelled/db"
)

type CMD string

const (
	CMD_HELLO CMD = "HELO"
	CMD_CONF  CMD = "CONF"
	CMD_PING  CMD = "PING"
	CMD_PONG  CMD = "PONG"
	CMD_DATA  CMD = "DATA"
	CMD_DGET  CMD = "DGET"
)

func BankToDbType(b string) db.BankType {
	switch b {
	case "TADC":
		return db.BankTypeADC
	case "COCL":
		return db.BankTypeContactClosure
	case "RELY":
		return db.BankTypeRelay
	case "TDAC":
		return db.BankTypeDAC
	}
	return db.BankTypeUnknown
}

type IncomeMessage struct {
	CMD  CMD
	SN   string
	IP   string
	Port uint16
	Body []rune
}

type OutcomeMessage struct {
	CMD  CMD
	IP   string
	Port uint16
	Body []byte
}

//
//func parsePingIncome(b []byte, n int) (*IncomeMessage, error) {
//	return &IncomeMessage{
//		Command: Command{
//			CMD: CMD(b[0:4]),
//		},
//		SN: string(b[4:12]),
//	}, nil
//}
