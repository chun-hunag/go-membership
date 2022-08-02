package postgres

import (
	"fmt"
	"go-membership/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type connection struct {
	configName string
	dsn        string
	db         *gorm.DB
}

func NewConnection(configName string) connection {
	return connection{configName: configName}
}

func (c *connection) open() error {
	if c.dsn == "" {
		connection := configs.GetDatabaseConnection(c.configName)
		c.dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Taipei", connection.Host, connection.Username, connection.Password, connection.Database, connection.Port)
	}
	db, err := gorm.Open(postgres.Open(c.dsn), &gorm.Config{})
	if err != nil {
		log.Println("Error Open Database Connection: postgres")
		return err
	}
	c.db = db
	return nil
}

func (c *connection) close() {
	if c.db == nil {
	}

	db, err := c.db.DB()
	if err != nil {
		log.Println("get sql.DB failed when close connection: postgres")
	}
	err = db.Close()
	if err != nil {
		log.Println("Failed when close connection: postgres")
	}
}
