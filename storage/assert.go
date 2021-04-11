package storage

type Assert struct {
	Amount  int64
	Cost    int64
	Name    string
	ValidTo string
}

func (s *SqLite) CreateAssert(assert *Assert) error {
	_, err := s.db.Exec(`INSERT INTO assert (name, amount, cost, valid_to) VALUES ($1, $2, $3, $4)`,
		assert.Name,
		assert.Amount,
		assert.Cost,
		assert.ValidTo)
	return err
}
