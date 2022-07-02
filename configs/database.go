package configs

import (
	"os"
)

type Connection struct {
	Driver     string
	Port       string
	Host       string
	Database   string
	Username   string
	Password   string
	Charset    string
	Collection string
}

func newConnection(driver, port, host, database, username, password string) Connection {
	return Connection{
		Driver:   driver,
		Port:     port,
		Host:     host,
		Database: database,
		Username: username,
		Password: password,
	}
}

var connections = make(map[string]Connection)

// loadDatabaseConfig Add all databases setting here
func loadDatabaseConfig() {
	connections["postgres"] = newConnection(
		"pgsql",
		os.Getenv("POSTGRESQL_PORT"),
		os.Getenv("POSTGRESQL_HOST"),
		os.Getenv("POSTGRESQL_DB"),
		os.Getenv("POSTGRESQL_USER"), os.Getenv("POSTGRESQL_PASSWORD"),
	)
}

func GetDatabaseConnection(name string) Connection {
	var connection Connection
	if connection, ok := connections[name]; ok {
		return connection
	}
	return connection
}
