package drift_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"github.com/Di-Argus/Drift"
	"github.com/Di-Argus/Drift/pkg/config"
)


type MockDriver struct {
	InitliazeMigrationsFunc  func() error
	GetAppliedMigrationsFunc func() (map[int64]string, error)
	ApplyFunc                func(version int64, name, checksum, sql string) error
	DownFunc                 func(version int64, sql string) error
	CloseFunc                func()
}

func (md *MockDriver) InitializeMigrations() error {
	if md.InitliazeMigrationsFunc != nil {
		return md.InitliazeMigrationsFunc()
	}
	return nil
}

func (md *MockDriver) GetAppliedMigrations() (map[int64]string, error) {
	if md.GetAppliedMigrationsFunc != nil {
		return md.GetAppliedMigrationsFunc()
	}
	version := "2006010215040"
	migrations := make(map[int64]string)
	for i := 0; i < 5; i++ {
		v, _ := strconv.Atoi(version)
		migrations[int64(v*10+i)] = "fake_checksum"
	}

	return migrations, nil
}

func (md *MockDriver) Apply(version int64, name, checksum, sql string) error {
	if md.ApplyFunc != nil {
		return md.ApplyFunc(version, name, checksum, sql)
	}
	if version != 20260301000000 {
		return fmt.Errorf("expected version 20260301000000 got %v", version)
	}
	return nil
}

func (md *MockDriver) Down(version int64, sql string) error {
	if md.DownFunc != nil {
		return md.DownFunc(version, sql)
	}
	return nil
}

func (md *MockDriver) Close() {
	if md.CloseFunc != nil {
		md.CloseFunc()
	}
	return
}

func TestRunUp(t *testing.T) {
	tmpDir := t.TempDir()

	migrationFile := filepath.Join(tmpDir, "20260301000000_create_users_table.up.sql")
	err := os.WriteFile(migrationFile, []byte("CREATE TABLE users (id INT);"), 0644)
	if err != nil {
		t.Fatalf("failed to create temp migration file: %v", err)
	}

	cfg := config.Config{
		Dir: tmpDir,
	}

	mockDrv := &MockDriver{}

	migrator := drift.Migrator{
		Config: &cfg,
		Driver: mockDrv,
	}

	err = migrator.RunUp()
	if err != nil {
		t.Errorf("unexpected error : %v", err)
	}
}
