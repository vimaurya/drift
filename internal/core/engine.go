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
	"golang.org/x/tools/go/cfg"
)


func RunUp() error {
	
	data, err := os.ReadFile(".gomigrate.json")
	if err!=nil{
		return fmt.Errorf("failed to read config : %w", err)
	}
	var cfg config.Config

	err = json.Unmarshal(data, &cfg)
	if err!=nil{
		return fmt.Errorf("failed to Unmarshal : %w", err)
	}

	appliedMigrations, err := d.GetAppliedMigrations()
	if err!=nil{
		return fmt.Errorf("failed to fetch applied migrations : %w", err)
	}

	availableMigrations, err := migrations.GetAvailableMigrations(migrationDir)
	if err!=nil{
		return fmt.Errorf("failed to fetch available migrations : %w", err)
	}

	for _, fileName := range availableMigrations {
		parts := strings.SplitN(fileName, "_", 2)
		
		version, err:= strconv.ParseInt(parts[0], 10, 64)
		if err!=nil{
			return fmt.Errorf("failed to read version of the migration %s : %w", fileName, err)
		}
		name := strings.TrimSuffix(parts[1], ".up.sql")

		sqlContent, err := os.ReadFile(filepath.Join(migrationDir, fileName))

    currentChecksum := calculateCheckSum(string(sqlContent))

		checksum, exists := appliedMigrations[version]

		if !exists{
			err = d.Apply(version, name, currentChecksum, string(sqlContent))
			if err!=nil{
				return fmt.Errorf("failed to apply migration %s : %w", fileName, err)
			}
		} else if checksum != currentChecksum{
    	return fmt.Errorf("applied migration file edited %s", fileName)
		}
	}

	return nil
}
