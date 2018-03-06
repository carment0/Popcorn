// Copyright (c) 2018 Popcorn
// Author(s) Calvin Feng, Carmen To

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

const DIR = "datasets/production/"

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {
	var dbCredentials string
	if os.Getenv("HEROKU_POSTGRESQL_BROWN_URL") != "" {
		dbCredentials = os.Getenv("HEROKU_POSTGRESQL_BROWN_URL")
	} else if os.Getenv("DATABASE_URL") != "" {
		dbCredentials = os.Getenv("DATABASE_URL")
	} else {
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
	var moviePopularityMap map[uint]map[string]float64
	var featuresMap map[uint][]float64
	var metadataMap map[uint]map[string]string
	var movieClusterMap map[uint]uint
	var movieClusterRelationMap map[uint]map[string][]string
	var loadError error

	movieModelsMap, loadError = loadMoviesCSVFile(DIR + "movies.csv")
	if loadError != nil {
		logrus.Fatal("Failed to load movie models from CSV data:", loadError)
		return
	} else {
		logrus.Info("Movie models are loaded from csv files")
	}

	moviePopularityMap, loadError = loadPopularityCSVFile(DIR + "popularity.csv")
	if loadError != nil {
		logrus.Fatal("Failed to load movie popularities from CSV data:", loadError)
		return
	} else {
		logrus.Info("Movie popularities are loaded from csv files")
	}

	metadataMap, loadError = loadMetadataCSVFile(DIR + "links.csv")
	if loadError != nil {
		logrus.Fatal("Failed to load movie metadata from CSV data:", loadError)
		return
	} else {
		logrus.Info("Movie metadata are loaded from csv files")
	}

	featuresMap, loadError = loadFeatureCSVFile(DIR + "features.csv")
	if loadError != nil {
		logrus.Error("Failed to load movie features from CSV data:", loadError)
	} else {
		logrus.Info("Movie features are loaded from csv files")
	}

	movieClusterMap, loadError = loadMovieClusterCSVFile(DIR + "clusters.csv")
	if loadError != nil {
		logrus.Error("Failed to load movie clusters from CSV data:", loadError)
	} else {
		logrus.Info("Movie clusters are loaded from csv files")
	}

	movieClusterRelationMap, loadError = loadMovieClusterRelationsCSVFile(DIR + "clusters.csv")
	if loadError != nil {
		logrus.Error("Failed to load movie clusters relations from CSV data:", loadError)
	} else {
		logrus.Info("Movie clusters relations are loaded from csv files")
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

		if dict, ok := moviePopularityMap[movieID]; ok {
			movie.AverageRating = dict["avg_rating"]
			movie.NumRating = int(dict["num_rating"])
		}

		if dict, ok := metadataMap[movieID]; ok {
			movie.IMDBID = dict["imdb"]
			movie.TMDBID = dict["tmdb"]
		}

		if movieClusterMap != nil {
			if value, ok := movieClusterMap[movieID]; ok {
				movie.ClusterID = value
			}
		}

		if featuresMap != nil {
			if value, ok := featuresMap[movieID]; ok {
				movie.Feature = value
			}
		}

		if dict, ok := movieClusterRelationMap[movieID]; ok {
			movie.NearestClusters = dict["closest"]
			movie.FarthestClusters = dict["farthest"]
		}

		if db.Create(movie).Error == nil {
			count += 1
		}
	}

	logrus.Infof("Completed seeding %d movies", count)
}
