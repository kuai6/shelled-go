package sqlite3

import (
	"github.com/jmoiron/sqlx"

	"shelled-backend/shelled/db"
)

type BankStorage struct {
	db *sqlx.DB
}

func NewBankStorage(db *sqlx.DB) *BankStorage {
	return &BankStorage{db: db}
}

func (s *BankStorage) FindByDeviceIdAndNumber(deviceId int32, number int32) (*db.Bank, error) {
	return nil, nil
}
