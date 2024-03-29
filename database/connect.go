package database

import (
	"fmt"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
}

func Connect(config *Config) {
	var (
		err     error
		port, _ = strconv.ParseUint(config.Port, 10, 32)
		dsn     = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.Host, port, config.User, config.Password, config.DBName,
		)
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		fmt.Println(
			err.Error(),
		)
		panic("failed to connect database")
	}

	// RunAutoMigrations()

	fmt.Println("Connection Opened to Database")
}

var DB *gorm.DB
