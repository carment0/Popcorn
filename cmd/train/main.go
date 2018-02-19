// Copyright (c) 2018 Popcorn
// Author(s) Calvin Feng

package main

import (
	"github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/mat"
	"popcorn/lowrank"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	processor, err := lowrank.NewDataProcessor("fullset/ratings.csv", "fullset/movies.csv")
	if err != nil {
		logrus.Fatal(err)
	}

	featureDim := 10
	R := processor.GetRatingMatrix()
	approx := lowrank.NewApproximator(R, featureDim)
	approx.DataProcessor = processor

	// Start training
	approx.Train(300, 1, 0, 1e-5)

	J, _ := approx.MovieLatent.Dims()
	featureMapByMovieID := make(map[int][]float64)
	for j := 0; j < J; j += 1 {
		movieID := processor.MovieIndexToID[j]
		features := make([]float64, featureDim)
		mat.Row(features, j, approx.MovieLatent)
		featureMapByMovieID[movieID] = features
	}

	writeFeaturesToCSV("fullset/features.csv", featureMapByMovieID, featureDim)
	writePopularityToCSV("fullset/popularity.csv", processor.MovieMap)
}
