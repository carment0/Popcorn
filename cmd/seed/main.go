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

	var movieModelsMap map[uint]*model.Movie
	var movieRatingsMap map[uint][]float64
	var featuresMap map[uint][]float64
	var metadataMap map[uint]map[string]string
	var loadError error

	movieModelsMap, loadError = LoadMoviesCSVFile("dataset/movies.csv")
	if loadError != nil {
		logrus.Fatal("Failed to load movie models from CSV data:", loadError)
		return
	} else {
		logrus.Info("Movie models are loaded from csv files")
	}

	movieRatingsMap, loadError = LoadRatingsCSVFile("dataset/ratings.csv")
	if loadError != nil {
		logrus.Fatal("Failed to load movie ratings from CSV data:", loadError)
		return
	} else {
		logrus.Info("Movie ratings map by movie ID are loaded from csv files")
	}

	metadataMap, loadError = LoadMetadataCSVFile("dataset/links.csv")
	if loadError != nil {
		logrus.Fatal("Failed to load movie metadata from CSV data:", loadError)
		return
	} else {
		logrus.Info("Movie metadata are loaded from csv files")
	}

	featuresMap, loadError = LoadFeatureCSVFile("dataset/features.csv")
	if loadError != nil {
		logrus.Error("Failed to load movie features from CSV data:", loadError)
	} else {
		logrus.Info("Movie features are loaded from csv files")
	}

	if db.HasTable(&model.Movie{}) {
		db.DropTable(&model.Movie{})
		logrus.Info("Existing \"movies\" table is dropped")
	}

	db.CreateTable(&model.Movie{})
	logrus.Info("New \"movies\" table is created")

	count := 0
	for movieID := range movieModelsMap {
		movie := movieModelsMap[movieID]

		if ratingList, ok := movieRatingsMap[movieID]; ok {
			movie.AverageRating = Average(ratingList)
			movie.NumRating = len(ratingList)
		}

		if dict, ok := metadataMap[movieID]; ok {
			movie.IMDBID = dict["imdb"]
			movie.TMDBID = dict["tmdb"]
		}

		if featuresMap != nil {
			if value, ok := featuresMap[movieID]; ok {
				movie.Feature = value
			}
		}

		if db.Create(movie).Error == nil {
			count += 1
		}
	}

	logrus.Infof("Completed seeding %d movies", count)

}
