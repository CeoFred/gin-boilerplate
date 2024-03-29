package database

import (
	"context"
	"gorm.io/gorm"
	"log"
	"time"
)

const (
	dbTimeout = 30 * time.Second
)

func RunManualMigration(db *gorm.DB) {

	query1 := `CREATE TABLE IF NOT EXISTS users (
			id UUID NOT NULL,
			email VARCHAR(255) NOT NULL,
			password VARCHAR(255),
			first_name VARCHAR(255),
			last_name VARCHAR(255),
			avatar_url VARCHAR(255) DEFAULT NULL,
			ip VARCHAR(255) DEFAULT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			last_login VARCHAR(255) NULL,
		 role VARCHAR(255) DEFAULT 'user',
		 email_verified BOOLEAN DEFAULT FALSE,
		 country VARCHAR(255) DEFAULT NULL,
		 phone_number VARCHAR(255) DEFAULT NULL,
		 status VARCHAR(255) DEFAULT 'Inactive'
			);
			`

	migrationQueries := []string{
		query1,
	}

	log.Println("running db migration :::::::::::::")

	_, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	for _, query := range migrationQueries {
		err := db.Exec(query).Error
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("complete db migration")
}
