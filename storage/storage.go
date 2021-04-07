package storage

type Storage interface {
	Close() error
	RunMigrations() error
}
