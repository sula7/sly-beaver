package storage

import (
	"context"
	"fmt"
	"time"
)

func (s *SqLite) RunMigrations() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancelFunc()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin ctx transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS user(login text, password text);`)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return fmt.Errorf("tx rollback: %w", err)
		}
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO user (login, password) VALUES ('admin', '1234'), ('guest', '4321')`)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return fmt.Errorf("tx rollback: %w", err)
		}
		return err
	}

	return tx.Commit()
}
