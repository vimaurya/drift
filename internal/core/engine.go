package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vimaurya/gomigrate/internal/config"
	"github.com/vimaurya/gomigrate/internal/driver"
	"github.com/vimaurya/gomigrate/migrations"
)

func RunUp() error {
	data, err := os.ReadFile(".gomigrate.json")
	if err != nil {
		return fmt.Errorf("failed to read config : %w", err)
	}
	var cfg config.Config

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return fmt.Errorf("failed to Unmarshal : %w", err)
	}

	d, err := driver.GetDriver(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to get driver : %w", err)
	}

	appliedMigrations, err := d.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to fetch applied migrations : %w", err)
	}

	availableMigrations, err := migrations.GetAvailableMigrations(cfg.Dir)
	if err != nil {
		return fmt.Errorf("failed to fetch available migrations : %w", err)
	}

	for _, fileName := range availableMigrations {
		parts := strings.SplitN(fileName, "_", 2)

		version, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return fmt.Errorf("failed to read version of the migration %s : %w", fileName, err)
		}
		name := strings.TrimSuffix(parts[1], ".up.sql")

		sqlContent, err := os.ReadFile(filepath.Join(cfg.Dir, fileName))
		if err != nil {
			return fmt.Errorf("could not read migration(s) : %s", fileName)
		}

		currentChecksum := calculateCheckSum(string(sqlContent))

		checksum, exists := appliedMigrations[version]

		if !exists {
			err = d.Apply(version, name, currentChecksum, string(sqlContent))
			if err != nil {
				return fmt.Errorf("failed to apply migration %s : %w", fileName, err)
			}
		} else if checksum != currentChecksum {
			return fmt.Errorf("applied migration file edited %s", fileName)
		}
	}

	return nil
}

func RunDown() error {
	data, err := os.ReadFile(".gomigrate.json")
	if err != nil {
		return fmt.Errorf("failed to read config : %w", err)
	}
	var cfg config.Config

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return fmt.Errorf("failed to Unmarshal : %w", err)
	}

	d, err := driver.GetDriver(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to get driver : %w", err)
	}

	appliedMigrations, err := d.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to fetch applied migrations : %w", err)
	}

	availableMigrations, err := migrations.GetAvailableDownMigrations(cfg.Dir)
	if err != nil {
		return fmt.Errorf("failed to fetch available migrations : %w", err)
	}

	if len(availableMigrations) == 0 {
		fmt.Println("no migrations to rollback")
		return nil
	}

	for _, fileName := range availableMigrations {
		parts := strings.SplitN(fileName, "_", 2)

		version, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return fmt.Errorf("failed to read version of the migration %s : %w", fileName, err)
		}

		if checksum, exists := appliedMigrations[version]; exists {
			upFileName := strings.Replace(fileName, ".down.sql", ".up.sql", 1)
			upContent, err := os.ReadFile(filepath.Join(cfg.Dir, upFileName))
			if err != nil {
				return fmt.Errorf("failed to check up migration for version : %d", version)
			}
			if calculateCheckSum(string(upContent)) != checksum {
				return fmt.Errorf("integrity error : upfile %s was edited. can not safely rollback", upFileName)
			}

			downContent, err := os.ReadFile(filepath.Join(cfg.Dir, fileName))
			if err != nil {
				return fmt.Errorf("failed to read file %s : %w", fileName, err)
			}
			err = d.Down(version, string(downContent))
			if err != nil {
				return err
			}

			fmt.Printf("rolled back : %s\n", fileName)
			return nil
		}
	}
	return nil
}
