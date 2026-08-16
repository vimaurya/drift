package main

import (
	"fmt"
	_ "github.com/Di-Argus/Drift/pkg/driver/mysql"
	_ "github.com/Di-Argus/Drift/pkg/driver/postgres"
	"log"
	"os"
)

// postgres://postgres:root@localhost:5432/test_db?sslmode=disable
// mysql://root:vikash@tcp(127.0.0.1:3306)/test_db
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: drift <command> [options]")
		fmt.Println("Commands: init, create, up, down")
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]
	var err error
	switch command {
	case "init":
		err = runInit(args)
	case "create":
		err = runCreate(args)
	case "up":
		err = runUp(args)
	case "down":
		err = runDown(args)
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		log.Fatal(err)
	}
}
