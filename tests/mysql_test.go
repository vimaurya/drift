package tests

import (
	"context"
	"testing"

	"github.com/Di-Argus/Drift/internal/driver"
	_ "github.com/Di-Argus/Drift/internal/driver/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestPostgresDriver_ExecutionWithTestcontainers(t *testing.T) {
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

	err = d.InitializeMigrations()
	require.NoError(t, err)

	version := int64(20260301000000)
	name := "create_users"
	checksum := "abc123checksum"
	upSQL := "CREATE TABLE users (id INT);"

	err = d.Apply(version, name, checksum, upSQL)
	assert.NoError(t, err, "Apply should execute successfully")

	applied, err := d.GetAppliedMigrations()
	assert.NoError(t, err)
	assert.Contains(t, applied, version)
	assert.Equal(t, checksum, applied[version])

	downSQL := "DROP TABLE users;"
	err = d.Down(version, downSQL)
	assert.NoError(t, err, "Down should execute successfully")

	appliedAfterDown, err := d.GetAppliedMigrations()
	assert.NoError(t, err)
	assert.NotContains(t, appliedAfterDown, version)
}
