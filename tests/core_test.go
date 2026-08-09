package tests

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/Di-Argus/Drift/internal/config"
	"github.com/Di-Argus/Drift/internal/core"
	"github.com/stretchr/testify/assert"
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

func TestRunUp_Success(t *testing.T) {
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

	err = core.RunUp(cfg, mockDrv)
	if err != nil {
		t.Errorf("unexpected error : %v", err)
	}
}

func TestRunUp_IntegrityFailure(t *testing.T) {
	tmpDir := t.TempDir()

	migrationFile := filepath.Join(tmpDir, "20260301000000_create_users_table.up.sql")
	err := os.WriteFile(migrationFile, []byte("CREATE TABLE users (id INT);"), 0644)
	if err != nil {
		t.Fatalf("failed to create temp migration file: %v", err)
	}

	cfg := config.Config{
		Dir: tmpDir,
	}

	mockDrv := &MockDriver{
		GetAppliedMigrationsFunc: func() (map[int64]string, error) {
			return map[int64]string{20260301000000: "corrupted_checksum"}, nil
		},
	}

	err = core.RunUp(cfg, mockDrv)
	assert.Error(t, err, "expected error due to checksum mismatch")
}

func TestRunUp_AlreadyApplied(t *testing.T) {
	tmpDir := t.TempDir()

	migrationFile := filepath.Join(tmpDir, "20260301000000_create_users_table.up.sql")
	err := os.WriteFile(migrationFile, []byte("CREATE TABLE users (id INT);"), 0644)
	if err != nil {
		t.Fatalf("failed to create temp migration file: %v", err)
	}

	cfg := config.Config{
		Dir: tmpDir,
	}

	applyCalled := false

	mockDrv := &MockDriver{
		GetAppliedMigrationsFunc: func() (map[int64]string, error) {
			sqlContent, _ := os.ReadFile(migrationFile)
			checksum := sha256.Sum256([]byte(sqlContent))

			return map[int64]string{20260301000000: hex.EncodeToString(checksum[:])}, nil
		},
		ApplyFunc: func(version int64, name, checksum, sql string) error {
			applyCalled = true
			return nil
		},
	}

	err = core.RunUp(cfg, mockDrv)
	assert.NoError(t, err)
	assert.False(t, applyCalled, "expected to not apply the migration due to already being applied")

}

func TestRunUp_FetchError(t *testing.T) {
	cfg := config.Config{Dir: t.TempDir()}

	mockDrv := &MockDriver{
		GetAppliedMigrationsFunc: func() (map[int64]string, error) {
			return nil, fmt.Errorf("connection timeout")
		},
	}

	err := core.RunUp(cfg, mockDrv)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection timeout")
}

func TestRunUp_ApplyFailure(t *testing.T) {
	tmpDir := t.TempDir()
	migrationFile := filepath.Join(tmpDir, "20260301000000_create_users_table.up.sql")
	err := os.WriteFile(migrationFile, []byte("INVALID SQL SYNTAX;"), 0644)
	assert.NoError(t, err)

	cfg := config.Config{Dir: tmpDir}

	mockDrv := &MockDriver{
		GetAppliedMigrationsFunc: func() (map[int64]string, error) {
			return map[int64]string{}, nil
		},
		ApplyFunc: func(version int64, name, checksum, sql string) error {
			return fmt.Errorf("syntax error near INVALID")
		},
	}

	err = core.RunUp(cfg, mockDrv)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "syntax error")
}

func TestRunDown_Success(t *testing.T) {
	tmpDir := t.TempDir()

	upFile := filepath.Join(tmpDir, "20260301000000_create_users_table.up.sql")
	downFile := filepath.Join(tmpDir, "20260301000000_create_users_table.down.sql")
	
	upContent := []byte("CREATE TABLE users (id INT);")
	downContent := []byte("DROP TABLE users;")

	assert.NoError(t, os.WriteFile(upFile, upContent, 0644))
	assert.NoError(t, os.WriteFile(downFile, downContent, 0644))

	cfg := config.Config{Dir: tmpDir}
	
	downCalled := false 

	mockDrv := &MockDriver{
		GetAppliedMigrationsFunc: func() (map[int64]string, error) {
			return map[int64]string{
				20260301000000: core.CalculateCheckSum(string(upContent)),
			}, nil
		},

		DownFunc: func(version int64, sql string) error {
			downCalled = true
			assert.Equal(t, int64(20260301000000), version)
			assert.Equal(t, "DROP TABLE users;", sql)
			return nil
		},	
	}

	err := core.RunDown(cfg, mockDrv)
	assert.NoError(t, err)
	assert.True(t, downCalled, "expected down driver method to be called")
}

func TestRunDown_NoMigrations(t *testing.T) {
	tmpDir := t.TempDir() 
	cfg := config.Config{Dir: tmpDir}

	mockDrv := &MockDriver{
		GetAppliedMigrationsFunc: func() (map[int64]string, error) {
			return map[int64]string{20260301000000:"some_checksum"}, nil
		},
	}

	err := core.RunDown(cfg, mockDrv)
	assert.NoError(t, err, "expected no error as there are no migration files")
}

func TestRunDown_IntegrityError(t *testing.T) {
	tmpDir := t.TempDir()

	upFile := filepath.Join(tmpDir, "20260301000000_create_users_table.up.sql")
	downFile := filepath.Join(tmpDir, "20260301000000_create_users_table.down.sql")

	assert.NoError(t, os.WriteFile(upFile, []byte("CREATE TABLE users (id INT);"), 0644))
	assert.NoError(t, os.WriteFile(downFile, []byte("DROP TABLE users;"), 0644))

	cfg := config.Config{Dir: tmpDir}

	mockDrv := &MockDriver{
		GetAppliedMigrationsFunc: func() (map[int64]string, error) {
			return map[int64]string{20260301000000: "tampered_checksum_value"}, nil
		},
	}

	err := core.RunDown(cfg, mockDrv)
	assert.Error(t, err, "expected an integrity error because the up migration file was modified")
	assert.Contains(t, err.Error(), "integrity error")
}
