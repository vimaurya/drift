package driver

type MigrationRecord struct {
	Version int64
	Checksum string
}

type Driver interface {
	Init() error
	GetAppliedMigrations() (map[int64]string, error)
	Apply(version int64, name, checksum, sql string) error
	Close() 
}
