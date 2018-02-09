// Copyright (c) 2018 Popcorn
// Author(s) Calvin Feng

package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"popcorn/model"
)

const (
	LocalDBUser     = "popcorn"
	LocalDBPassword = "popcorn"
	LocalDBName     = "popcorn_development"
	LocalSSLMode    = "disable"
)

// SetupDatabase will perform database connection and auto migration on all gorm.Models
func SetupDatabase() (*gorm.DB, error) {
	dbCredentials := os.Getenv("DATABASE_URL")

	if dbCredentials == "" {
		dbCredentials = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
			LocalDBUser, LocalDBPassword, LocalDBName, LocalSSLMode,
		)
	}

	db, err := gorm.Open("postgres", dbCredentials)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.Movie{}, &model.User{}, &model.Rating{})

	return db, nil
}
