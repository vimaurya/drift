package driver

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type PostgresDriver struct {
	db *sql.DB
}

func NewPostgresDriver(url string) (*PostgresDriver, error) {
	db, err := sql.Open("postgres", url)
	if err!=nil{
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err!=nil{
		return nil, err
	}
	
	return &PostgresDriver{
		db: db,
	}, nil
}

func (p *PostgresDriver) Init() error{
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

// To-do 
// Implement GetAppliedMigrations and Apply

func (p *PostgresDriver) GetAppliedMigrations() (map[int64]string, error) {
	var migrationRecord = make(map[int64]string)
	
	tx, err := p.db.Begin()

	query := `
	SELECT version, checksum from schema_migrations order by version ASC;
	`
  
	rows, err := tx.Query(query)
	if err!=nil{
		return nil, err
	}
	
	defer rows.Close()

	for rows.Next() {
		var record MigrationRecord
		err := rows.Scan(
			&record.Version,
			&record.Checksum,
			)

		if err!=nil{
			return nil, err
		}

		migrationRecord[record.Version] = record.Checksum
	}

	return migrationRecord, nil
}

func (p *PostgresDriver) Apply(version int64, name, checksum, sqlContent string) error{
	tx, err := p.db.Begin()
	if err!=nil{
		return err
	}
	
	if _, err := tx.Exec(sqlContent); err!=nil{
		tx.Rollback()
		return fmt.Errorf("failed to execute migration %s: %w", name, err)
	}

	query := `
		INSERT INTO schema_migrations (version, name, checksum)
		VALUES ($1, $2, $3)
	`

	if _, err := tx.Exec(query, version, name, checksum); err!=nil{
		tx.Rollback()
		return fmt.Errorf("failed to log migration to schema_migrations : %w", err)
	}

	return tx.Commit()
}

	

