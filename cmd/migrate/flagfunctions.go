package main

import (
	"flag"
	"fmt"
	"github.com/di-argus/drift/pkg/config"
	"github.com/di-argus/drift/pkg/core"
	"github.com/di-argus/drift/pkg/driver"
	"github.com/di-argus/drift/pkg/migration"
)

func runInit(args []string) error {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)

	initURL := initCmd.String("url", "", "Database URL")
	pathFlag := initCmd.String("path", "DB_Migrations", "Directory where migrations are saved")

	initCmd.Parse(args)
	if *initURL == "" {
		return fmt.Errorf("database url is required")
	}

	cfg := config.Config{
		DatabaseURL: *initURL,
		Dir:         *pathFlag,
	}

	err := config.Save(cfg)
	if err != nil {
		return fmt.Errorf("failed to save config : %v", err)
	}

	nDriver, err := driver.GetDriver(*initURL)
	if err != nil {
		return fmt.Errorf("failed to fetch driver : %v", err)
	}
	defer nDriver.Close()

	err = nDriver.InitializeMigrations()
	if err != nil {
		return fmt.Errorf("failed to init table : %v", err)
	}

	fmt.Println("Database initialized successfully.")

	return nil
}

func runCreate(args []string) error {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	nameFlag := createCmd.String("name", "", "Name of the migration")

	createCmd.Parse(args)

	if *nameFlag == "" {
		return fmt.Errorf("name of the migration can not be empty")
	}

	err := migration.Create(*nameFlag)
	if err != nil {
		return fmt.Errorf("failed to create the migration file")
	}
	fmt.Println("successfully created the migration files.")

	return nil
}

func runUp(args []string) error {
	upCmd := flag.NewFlagSet("up", flag.ExitOnError)
	upCmd.Parse(args)

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config : %w", err)
	}

	d, err := driver.GetDriver(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to get driver : %w", err)
	}
	err = core.RunUp(cfg, d)
	if err != nil {
		return fmt.Errorf("failed to make migration(s) : %v", err)
	}

	fmt.Println("successfully made all migartions")

	return nil
}

func runDown(args []string) error {
	downCmd := flag.NewFlagSet("down", flag.ExitOnError)
	downCmd.Parse(args)

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config : %w", err)
	}

	d, err := driver.GetDriver(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to get driver : %w", err)
	}

	err = core.RunDown(cfg, d)
	if err != nil {
		return fmt.Errorf("failed to down migration(s) : %v", err)
	}

	return nil
}

func printUsage() {
	fmt.Println("Usage: drift <command> [options]")
	fmt.Println("Commands:")
	fmt.Println("  init   - Initialize configuration and database")
	fmt.Println("  create - Create a new migration file")
	fmt.Println("  up     - Apply pending migrations")
	fmt.Println("  down   - Roll back the last migration")
}
