package main

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Missing argument 1")
	}

	action := args[1]
	loadEnv()
	host := os.Getenv("POSTGRESQL_HOST")
	port := os.Getenv("POSTGRESQL_PORT")
	user := os.Getenv("POSTGRESQL_USER")
	password := os.Getenv("POSTGRESQL_PASSWORD")
	database := os.Getenv("POSTGRESQL_DB")
	m, err := migrate.New(
		"file://database/migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s? sslmode=disable", user, password, host, port, database))
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

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
