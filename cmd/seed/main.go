// Copyright (c) 2018 Popcorn
// Author(s) Calvin Feng

package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	"os"
	"popcorn/model"
)

const (
	LocalDBUser     = "popcorn"
	LocalDBPassword = "popcorn"
	LocalDBName     = "popcorn_development"
	LocalSSLMode    = "disable"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	dbCredentials := os.Getenv("DATABASE_URL")

	if dbCredentials == "" {
		dbCredentials = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
			LocalDBUser, LocalDBPassword, LocalDBName, LocalSSLMode,
		)
	}

	db, err := gorm.Open("postgres", dbCredentials)
	if err != nil {
		logrus.Fatal("Cannot connect to database for seeding:", err)
	} else {
		logrus.Infof("Connected to database with the following credentials: %s", dbCredentials)
	}

	if db.HasTable(&model.Movie{}) {
		db.DropTable(&model.Movie{})
		logrus.Info("Dropped movies table")
	}

	db.CreateTable(&model.Movie{})
	logrus.Info("New movies table is created")

	if movies, err := LoadMovieModelsFromCSVData("dataset/movies.csv"); err != nil {
		fmt.Println("Fatal", err)
	} else {
		logrus.Info("Movies are loaded from csv files")

		count := 0
		for movieId := range movies {
			if db.Create(movies[movieId]).Error == nil {
				count += 1
			}
		}

		logrus.Infof("Completed seeding %d movies", count)
	}
}