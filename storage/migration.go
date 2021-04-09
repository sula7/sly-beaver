package storage

import (
	"context"
	"fmt"
	"time"
)

func (s *SqLite) RunMigrations() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	var isDBActual bool
	err := s.db.QueryRow(`SELECT exists(SELECT name FROM sqlite_master WHERE name = 'user')`).Scan(&isDBActual)
	if err != nil {
		return fmt.Errorf("check db version: %w", err)
	}

	if isDBActual {
		return nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin ctx transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS user(login text, password text, is_admin boolean);`)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return fmt.Errorf("tx rollback: %w", err)
		}
		return fmt.Errorf("create table user: %w", err)
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO user (login, password, is_admin)
			VALUES ('admin', '1234', true), ('guest', '4321', false)`)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return fmt.Errorf("tx rollback: %w", err)
		}
		return fmt.Errorf("insert users: %w", err)
	}

	return tx.Commit()
}
