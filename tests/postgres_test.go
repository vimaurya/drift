package tests

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Di-Argus/Drift/internal/driver"
	_ "github.com/Di-Argus/Drift/internal/driver/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestPostgresDriver_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)
	defer func() {
		_ = testcontainers.TerminateContainer(pgContainer)
	}()

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	d, err := driver.GetDriver(dsn)
	require.NoError(t, err)
	defer d.Close()
	assert.NotNil(t, d)

	err = d.InitializeMigrations()
	assert.NoError(t, err, "should initialize schema_migrations table successfully")

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	defer db.Close()

	var tableName string
	err = db.QueryRow("SELECT tablename FROM pg_tables WHERE tablename = 'schema_migrations'").Scan(&tableName)
	assert.NoError(t, err)
	assert.Equal(t, "schema_migrations", tableName)
}
