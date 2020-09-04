package sqlite3

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func InitStorage(db *sqlx.DB) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tx, err := db.BeginTxx(ctx, &sql.TxOptions{})
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()
	if err != nil {
		return fmt.Errorf("create transaction error: %s", err)
	}

	query := `
		CREATE TABLE IF NOT EXISTS device (
			id INTEGER PRIMARY KEY,
			name VARCHAR(255) DEFAULT NULL,
			ip   VARCHAR(255) NOT NULL,
			port INTEGER NOT NULL,
			serial_number CHAR (8) UNIQUE NOT NULL ,
			last_registered_at DATETIME DEFAULT NULL,
			last_ping_at DATETIME DEFAULT NULL
		);
	`
	if _, err := tx.Exec(query); err != nil {
		return fmt.Errorf("create schema error: %s", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction error: %s", err)
	}

	return nil
}
