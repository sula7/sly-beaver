package storage

type Storage interface {
	Close() error
	RunMigrations() error

	CheckPassword(login, password string) (isExists, isAdmin bool, err error)
	CreateAssert(assert *Assert) error
	GetNotDeletedAsserts() ([]*Assert, error)
	AddRemoveReason(assert *Assert) error
	GetAllRowsForCSV() ([]*Assert, error)
}
