package db

type BankType string

const (
	BankTypeADC            BankType = "adc"
	BankTypeContactClosure BankType = "contact_closure"
	BankTypeRelay          BankType = "relay"
	BankTypeDAC            BankType = "dac"
	BankTypeUnknown        BankType = "unknown"
)

type Bank struct {
	ID       int32    `db:"id"`
	Number   int      `db:"number"`
	Type     BankType `db:"type"`
	Pins     int      `db:"pins"`
	DeviceID int32    `db:"device_id"`
	State    []byte   `db:"state"`
}

type BankStorage interface {
	FindByDeviceIdAndNumber(deviceId int32, number int32) (*Bank, error)
}
