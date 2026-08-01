package driver

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/Di-Argus/Drift/internal/driver"
)

type PostgresDriver struct {
	db *sql.DB
}

func init(){
	driver.DriverRegister("postgres", New)
}

func New(db *sql.DB) (driver.Driver, error){
	return &PostgresDriver{db:db}, nil
}

func (p *PostgresDriver) InitializeMigrations() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version BIGINT PRIMARY KEY,
			name TEXT NOT NULL,
			checksum TEXT NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := p.db.Exec(query)
	return err
}

func (p *PostgresDriver) Close() {
	p.db.Close()
}

func (p *PostgresDriver) GetAppliedMigrations() (map[int64]string, error) {
	migrationRecord := make(map[int64]string)

	tx, err := p.db.Begin()

	query := `
	SELECT version, checksum from schema_migrations order by version ASC;
	`

	rows, err := tx.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var record driver.MigrationRecord
		err := rows.Scan(
			&record.Version,
			&record.Checksum,
		)
		if err != nil {
			return nil, err
		}

		migrationRecord[record.Version] = record.Checksum
	}

	return migrationRecord, nil
}

func (p *PostgresDriver) Apply(version int64, name, checksum, sqlContent string) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(sqlContent); err != nil {
		tx.Rollback()
		return err 	
	}

	query := `
		INSERT INTO schema_migrations (version, name, checksum)
		VALUES ($1, $2, $3)
	`

	if _, err := tx.Exec(query, version, name, checksum); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (p *PostgresDriver) Down(version int64, sqlContent string) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(sqlContent); err != nil {
		tx.Rollback()
		return err
	}

	query := `
		DELETE FROM schema_migrations WHERE version = $1
	`
	if _, err := tx.Exec(query, version); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
