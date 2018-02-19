// Copyright (c) 2018 Popcorn
// Author(s) Calvin Feng

package main

import (
	"encoding/csv"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"popcorn/lowrank"
	"strconv"
)

func writePopularityToCSV(filepath string, movieMap map[int]*lowrank.Movie) error {
	csvFile, fileErr := os.Create(filepath)
	if fileErr != nil {
		return fileErr
	}

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Write the header first
	var writerError error
	writerError = writer.Write([]string{"movieId", "avgRating", "numRating"})
	if writerError != nil {
		return writerError
	}

	for movieID, movie := range movieMap {
		row := []string{
			strconv.Itoa(movieID),
			strconv.FormatFloat(movie.AvgRating, 'f', 2, 64),
			strconv.Itoa(len(movie.Ratings)),
		}

		writerError = writer.Write(row)
		if writerError != nil {
			logrus.Errorf("Failed to write row: %s\n", row)
		}
	}

	return nil
}

func writeFeaturesToCSV(filepath string, movieFeatures map[int][]float64, featureDim int) error {
	csvFile, fileErr := os.Create(filepath)
	if fileErr != nil {
		return fileErr
	}

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	header := []string{"movieId"}
	for i := 1; i <= featureDim; i += 1 {
		header = append(header, fmt.Sprintf("f%v", i))
	}

	var writerError error

	// Write the header first
	writerError = writer.Write(header)
	if writerError != nil {
		return writerError
	}

	for movieID, features := range movieFeatures {
		row := []string{strconv.Itoa(movieID)}

		for _, feature := range features {
			row = append(row, strconv.FormatFloat(feature, 'f', 6, 64))
		}

		writerError = writer.Write(row)
		if writerError != nil {
			logrus.Errorf("Failed to write row: %s\n", row)
		}
	}

	return nil
}
