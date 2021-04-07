package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type SqLite struct {
	Storage
	db *sql.DB
}

func CreateDBFile(filePath string) error {
	_, err := os.Stat(filePath)
	if err == nil {
		return nil
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create .db file: %w", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("close .db file: %w", err)
	}

	return err
}

func (s *SqLite) Close() error {
	return s.db.Close()
}

func OpenDB(filePath string) (*SqLite, error) {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return nil, fmt.Errorf("open db conn: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return &SqLite{db: db}, err
}
