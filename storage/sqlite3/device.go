package sqlite3

import (
	"database/sql"
	"fmt"
	"time"

	//"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"shelled-backend/shelled/db"
)

type DeviceStorage struct {
	db *sqlx.DB
}

func NewDeviceStorage(db *sqlx.DB) *DeviceStorage {
	return &DeviceStorage{
		db: db,
	}
}

func (s *DeviceStorage) FindBySerialNumber(sn string) (*db.Device, error) {
	query := `SELECT * FROM device WHERE serial_number = $1`
	d := &db.Device{}
	err := s.db.Get(d, query, sn)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get device by sn error: %s", err)
	}

	return d, nil
}

func (s *DeviceStorage) Create(sn string, ip string, port uint16) (*db.Device, error) {
	d := &db.Device{}
	query := `
		INSERT INTO device 
			(serial_number, ip, port, last_registered_at)
			VALUES
			($1, $2, $3, $4);
	`
	result, err := s.db.Exec(query, sn, ip, port, time.Now())
	if err != nil {
		return nil, fmt.Errorf("insert error: %s", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("get last insert id error: %s", err)
	}

	if err := s.db.Get(d, `SELECT * FROM device WHERE id = $1`, id); err != nil {
		return nil, fmt.Errorf("get row error: %s", err)
	}

	return d, nil
}

func (s *DeviceStorage) Update(device *db.Device) (*db.Device, error) {
	query := `UPDATE 
    			device 
			  SET name = :name,
			      ip = :ip,
			      port = :port,
			      last_registered_at = :last_registered_at, 
			      last_ping_at = :last_ping_at 
			  WHERE id = :id`
	if _, err := s.db.NamedExec(query, device); err != nil {
		return nil, fmt.Errorf("update row error: %s", err)
	}
	return device, nil
}
