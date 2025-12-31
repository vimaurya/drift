package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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
	case "up":
		upCmd := flag.NewFlagSet("up", flag.ExitOnError)

		upURL := upCmd.String("url", "", "Database URL")

		upCmd.Parse(os.Args[2:])
		if *upURL == "" {
			log.Fatal("database url is required")
		}
		fmt.Println("fetching driver...")
		nDriver, err := driver.GetDriver(*upURL)
		if err != nil {
			log.Fatalf("failed to fetch driver : %v", err)
		}
		defer nDriver.Close()
		fmt.Println("fetched driver successfully")
		
		err = nDriver.Init()

		fmt.Println("initializing table..")
		if err != nil {
			log.Fatalf("failed to init table : %v", err)
		}

		fmt.Println("Database initialized successfully.")
	
	case "create":
		createCmd := flag.NewFlagSet("create", flag.ExitOnError)
		nameFlag := createCmd.String("name", "", "Name of the migration")
		pathFlag := createCmd.String("path", "", "Directory where migrations are saved")
		
		createCmd.Parse(os.Args[2:])

		if *nameFlag == ""{
			log.Fatal("name of the migration can not be empty")
		}

		err := migration.Create(*nameFlag, *pathFlag)
		if err!=nil{
			log.Fatal("failed to create the migration file")
		}
	
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])

	}
}
