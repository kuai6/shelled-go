package sqlite3

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"shelled-backend/shelled/db"
)

type BankStorage struct {
	db *sqlx.DB
}

func NewBankStorage(db *sqlx.DB) *BankStorage {
	return &BankStorage{db: db}
}

func (s *BankStorage) Register(deviceId int32, bankType db.BankType, number, pins int) (*db.Bank, error) {

	state := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	query := `	
		INSERT INTO 
			bank (device_id, number, type, pins, state)
		VALUES 
			($1, $2, $3, $4, $5);
	`

	result, err := s.db.Exec(query, deviceId, number, bankType, pins, state)
	if err != nil {
		return nil, fmt.Errorf("insert error: %s", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("get last insert id error: %s", err)
	}
	b := &db.Bank{}
	if err := s.db.Get(b, `SELECT * FROM bank WHERE id = $1`, id); err != nil {
		return nil, fmt.Errorf("get row error: %s", err)
	}

	return b, nil
}

func (s *BankStorage) FindByDeviceIdAndNumber(deviceId int32, number int) (*db.Bank, error) {
	query := `SELECT * FROM bank WHERE device_id = $1 and number = $2`
	b := &db.Bank{}
	err := s.db.Get(b, query, deviceId, number)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get bank by device id and number error: %s", err)
	}

	return b, nil
}

func (s *BankStorage) Update(bank *db.Bank) (*db.Bank, error) {
	query := `UPDATE 
    			bank 
			  SET state = :state
			  WHERE id = :id`
	if _, err := s.db.NamedExec(query, bank); err != nil {
		return nil, fmt.Errorf("update row error: %s", err)
	}
	return bank, nil
}
