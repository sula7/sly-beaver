package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type SqLite struct {
	Storage
	db *sql.DB
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

func (s *SqLite) CheckPassword(login, password string) (isExists, isAdmin bool, err error) {
	err = s.db.QueryRow(`SELECT exists(SELECT login
			FROM user
			WHERE login = $1
			  AND password = $2), is_admin
			  FROM user`, login, password).
		Scan(&isExists, &isAdmin)
	return
}
