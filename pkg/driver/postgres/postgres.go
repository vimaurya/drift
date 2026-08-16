package driver

import (
	"context"
	"database/sql"
	"github.com/Di-Argus/Drift/pkg/driver"
	_ "github.com/lib/pq"
	"time"
)

type PostgresDriver struct {
	db *sql.DB
}

func init() {
	driver.DriverRegister("postgres", New)
}

func New(connURL string) (driver.Driver, error) {
	db, err := sql.Open("postgres", connURL)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return &PostgresDriver{db: db}, nil
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
