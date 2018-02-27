// Copyright (c) 2018 Popcorn
// Author(s) Calvin Feng

package main

import (
	"github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/mat"
	"popcorn/lowrank"
)

const InputDir = "datasets/100k/"
const OutputDir = "datasets/100k/"

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {
	processor, err := lowrank.NewMatrixConverter(InputDir+"ratings.csv", InputDir+"movies.csv")
	if err != nil {
		logrus.Fatal(err)
	}

	featureDim := 10
	R := processor.GetRatingMatrix()
	fact := lowrank.NewFactorizer(R, featureDim)
	fact.MatrixConverter = processor

	// Start training
	fact.Train(100, 10, 0.03, 1e-5)

	J, _ := fact.MovieLatent.Dims()
	featureMapByMovieID := make(map[int][]float64)
	for j := 0; j < J; j += 1 {
		movieID := processor.MovieIndexToID[j]
		features := make([]float64, featureDim)
		mat.Row(features, j, fact.MovieLatent)
		featureMapByMovieID[movieID] = features
	}

	writeFeaturesToCSV(OutputDir+"features.csv", featureMapByMovieID, featureDim)
	writePopularityToCSV(OutputDir+"popularity.csv", processor.MovieMap)
}
