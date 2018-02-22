// Copyright (c) 2018 Popcorn
// Author(s) Calvin Feng

package main

import (
	"github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/mat"
	"popcorn/lowrank"
)

const INPUT_DIR = "datasets/26m/"
const OUTPUT_DIR = "datasets/26m/"

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {
	processor, err := lowrank.NewMatrixConverter(INPUT_DIR+"ratings.csv", INPUT_DIR+"movies.csv")
	if err != nil {
		logrus.Fatal(err)
	}

	featureDim := 10
	R := processor.GetRatingMatrix()
	approx := lowrank.NewFactorizer(R, featureDim)
	approx.MatrixConverter = processor

	// Start training
	approx.Train(400, 10, 0.03, 1e-5)

	J, _ := approx.MovieLatent.Dims()
	featureMapByMovieID := make(map[int][]float64)
	for j := 0; j < J; j += 1 {
		movieID := processor.MovieIndexToID[j]
		features := make([]float64, featureDim)
		mat.Row(features, j, approx.MovieLatent)
		featureMapByMovieID[movieID] = features
	}

	writeFeaturesToCSV(OUTPUT_DIR+"features.csv", featureMapByMovieID, featureDim)
	writePopularityToCSV(OUTPUT_DIR+"popularity.csv", processor.MovieMap)
}
