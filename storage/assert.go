package storage

type Assert struct {
	ID           int64
	Amount       int64
	Cost         int64
	Name         string
	ValidTo      string
	CreatedAt    string
	RemoveReason string
}

func (s *SqLite) CreateAssert(assert *Assert) error {
	_, err := s.db.Exec(`INSERT INTO assert (name, amount, cost, valid_to) VALUES ($1, $2, $3, $4)`,
		assert.Name,
		assert.Amount,
		assert.Cost,
		assert.ValidTo)
	return err
}

func (s *SqLite) GetNotDeletedAsserts() ([]*Assert, error) {
	rows, err := s.db.Query(`SELECT id, created_at, name, amount FROM assert WHERE remove_reason IS NULL`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	asserts := []*Assert{}
	for rows.Next() {
		a := Assert{}
		err = rows.Scan(&a.ID, &a.CreatedAt, &a.Name, &a.Amount)
		if err != nil {
			return nil, err
		}

		asserts = append(asserts, &a)
	}

	return asserts, err
}

func (s *SqLite) AddRemoveReason(assert *Assert) error {
	_, err := s.db.Exec(`UPDATE assert SET remove_reason = $1, removed_at = CURRENT_DATE WHERE id = $2`,
		assert.RemoveReason, assert.ID)
	return err
}

func (s *SqLite) GetLastWeekAllAsserts() ([]*Assert, error) {
	rows, err := s.db.Query(`SELECT name, amount, cost, valid_to
			FROM assert
			WHERE DATE(created_at) >= DATE('now', '-7 days')`)
	if err != nil {
		return nil, err
	}

	asserts := []*Assert{}

	defer rows.Close()

	for rows.Next() {
		a := &Assert{}
		err = rows.Scan(&a.Name, &a.Amount, &a.Cost, &a.ValidTo)
		if err != nil {
			return nil, err
		}

		asserts = append(asserts, a)
	}

	return asserts, err
}

func (s *SqLite) GetLastWeekRemovedAsserts() ([]*Assert, error) {
	rows, err := s.db.Query(`SELECT name, amount, remove_reason
			FROM assert
			WHERE DATE(removed_at) >= DATE('now', '-7 days')`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	asserts := []*Assert{}
	for rows.Next() {
		a := &Assert{}
		err = rows.Scan(&a.Name, &a.Amount, &a.RemoveReason)
		if err != nil {
			return nil, err
		}

		asserts = append(asserts, a)
	}

	return asserts, err
}
