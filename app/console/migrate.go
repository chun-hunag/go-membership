package main

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go-membership/app/utilis"
	"go-membership/configs"
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Missing argument 1")
	}

	initialize()
	action := args[1]
	connection := configs.GetDatabaseConnection("postgres")
	m, err := migrate.New(
		"file://database/migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s? sslmode=disable", connection.Username, connection.Password, connection.Host, connection.Port, connection.Database))
	if err != nil {
		log.Fatal("Migration instanced failed.")
	}

	switch action {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
		break
	case "down":
		if err := m.Down(); err != nil {
			log.Fatal(err)
		}
		break
	default:
		log.Fatal("Invalid argument: ", action, ". only support: up, down")
	}
}

func initialize() {
	utilis.LoadEnv()
	configs.LoadConfig()
}
