package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/vimaurya/gomigrate/internal/config"
	"github.com/vimaurya/gomigrate/internal/core"
	"github.com/vimaurya/gomigrate/internal/driver"
	"github.com/vimaurya/gomigrate/internal/migration"
)

// postgres://postgres:root@localhost:5432/test_db?sslmode=disable
func main() {
		if len(os.Args) < 2 {
		fmt.Println("Usage: migrate <command> [options]")
		fmt.Println("Commands: up, create")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		initCmd := flag.NewFlagSet("init", flag.ExitOnError)

		initURL := initCmd.String("url", "", "Database URL")
		pathFlag := initCmd.String("path", "DB_Migrations", "Directory where migrations are saved")

		initCmd.Parse(os.Args[2:])
		if *initURL == "" {
			log.Fatal("database url is required")
		}

		cfg := config.Config{
			DatabaseURL: *initURL,
      Dir: *pathFlag,
		}

		err := config.Save(cfg)
		if err!=nil{
			log.Fatalf("fialed to save config : %v", err)
		}

		nDriver, err := driver.GetDriver(*initURL)
		if err != nil {
			log.Fatalf("failed to fetch driver : %v", err)
		}
		defer nDriver.Close()
		
		err = nDriver.Init()

		if err != nil {
			log.Fatalf("failed to init table : %v", err)
		}

		fmt.Println("Database initialized successfully.")
	
	case "create":
		createCmd := flag.NewFlagSet("create", flag.ExitOnError)
		nameFlag := createCmd.String("name", "", "Name of the migration")
		
		createCmd.Parse(os.Args[2:])

		if *nameFlag == ""{
			log.Fatal("name of the migration can not be empty")
		}

		err := migration.Create(*nameFlag)
		if err!=nil{
			log.Fatal("failed to create the migration file")
		}
		fmt.Println("successfully created the migration files.")
		
	case "up":
		upCmd := flag.NewFlagSet("up", flag.ExitOnError)
		upCmd.Parse(os.Args[2:])
		err := core.RunUp()
		if err!=nil{
			log.Fatalf("failed to make migration(s) : %v", err)
		}
		
		fmt.Println("successfully made all migartions")

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])

	}
}
