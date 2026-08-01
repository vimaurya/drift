package driver

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySqlDriver struct {
	db *sql.DB
}

func NewMySQLDriver(url string) (*MySqlDriver, error) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, fmt.Errorf("failed to open mysql connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping mysql: %w", err)
	}

	return &MySqlDriver{db: db}, nil
}

func (m *MySqlDriver) Init() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version BIGINT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			checksum VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := m.db.Exec(query)
	return err
}

func (m *MySqlDriver) GetAppliedMigrations() (map[int64]string, error) {
	migrationRecord := make(map[int64]string)

	query := `SELECT version, checksum FROM schema_migrations ORDER BY version ASC`

	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var version int64
		var checksum string
		if err := rows.Scan(&version, &checksum); err != nil {
			return nil, err
		}
		migrationRecord[version] = checksum
	}

	return migrationRecord, nil
}

func (m *MySqlDriver) Apply(version int64, name, checksum, sqlContent string) error {
	if _, err := m.db.Exec(sqlContent); err != nil {
		return fmt.Errorf("migration execution failed: %w", err)
	}

	query := `INSERT INTO schema_migrations (version, name, checksum) VALUES (?, ?, ?)`
	_, err := m.db.Exec(query, version, name, checksum)
	if err != nil {
		return fmt.Errorf("failed to log migration version: %w", err)
	}

	return nil
}

func (m *MySqlDriver) Down(version int64, sqlContent string) error {
	if _, err := m.db.Exec(sqlContent); err != nil {
		return fmt.Errorf("rollback execution failed: %w", err)
	}

	query := `DELETE FROM schema_migrations WHERE version = ?`
	_, err := m.db.Exec(query, version)
	if err != nil {
		return fmt.Errorf("failed to delete migration log: %w", err)
	}

	return nil
}

func (m *MySqlDriver) Close() {
	m.db.Close()
}
