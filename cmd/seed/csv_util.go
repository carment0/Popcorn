// Copyright (c) 2018 Popcorn
// Author(s) Calvin Feng

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/lib/pq"
	"io"
	"os"
	"popcorn/model"
	"regexp"
	"strconv"
	"strings"
)

func LoadFeatureCSVFile(filepath string) (map[uint][]float64, error) {
	if csvFile, err := os.Open(filepath); err != nil {
		return nil, err
	} else {
		reader := csv.NewReader(bufio.NewReader(csvFile))
		featureByMovieID := make(map[uint][]float64)
		for {
			var rowRecord []string
			var readerErr error

			rowRecord, readerErr = reader.Read()
			if readerErr != nil {
				if readerErr == io.EOF {
					break
				} else {
					fmt.Printf("Unexpected reader error: %v\n", readerErr)
					continue
				}
			}

			var movieID int64
			var parseErr error

			movieID, parseErr = strconv.ParseInt(rowRecord[0], 10, 64)
			if parseErr != nil {
				continue
			}

			featureVector := []float64{}
			for i := 1; i < len(rowRecord); i += 1 {
				if value, err := strconv.ParseFloat(rowRecord[i], 64); err == nil {
					featureVector = append(featureVector, value)
				}
			}

			if len(featureVector) != len(rowRecord)-1 {
				continue
			}

			featureByMovieID[uint(movieID)] = featureVector
		}

		return featureByMovieID, nil
	}
}

func LoadMetadataCSVFile(filepath string) (map[uint]map[string]string, error) {
	if csvFile, err := os.Open(filepath); err != nil {
		return nil, err
	} else {
		reader := csv.NewReader(bufio.NewReader(csvFile))
		metadataByMovieID := make(map[uint]map[string]string)
		for {
			var rowRecord []string
			var readerErr error

			rowRecord, readerErr = reader.Read()
			if readerErr != nil {
				if readerErr == io.EOF {
					break
				} else {
					fmt.Printf("Unexpected reader error: %v\n", readerErr)
					continue
				}
			}

			var movieID int64
			var parseErr error

			movieID, parseErr = strconv.ParseInt(rowRecord[0], 10, 64)
			if parseErr != nil {
				continue
			}

			if _, ok := metadataByMovieID[uint(movieID)]; !ok {
				metadataByMovieID[uint(movieID)] = make(map[string]string)
			}

			metadataByMovieID[uint(movieID)]["imdb"] = "tt" + rowRecord[1]
			metadataByMovieID[uint(movieID)]["tmdb"] = rowRecord[2]
		}

		return metadataByMovieID, nil
	}
}

func LoadRatingsCSVFile(filepath string) (map[uint][]float64, error) {
	if csvFile, csvErr := os.Open(filepath); csvErr != nil {
		return nil, csvErr
	} else {
		reader := csv.NewReader(bufio.NewReader(csvFile))
		ratingsByMovieID := make(map[uint][]float64)
		for {
			var rowRecord []string
			var readerErr error

			rowRecord, readerErr = reader.Read()
			if readerErr != nil {
				if readerErr == io.EOF {
					break
				} else {
					fmt.Printf("Unexpected reader error: %v\n", readerErr)
					continue
				}
			}

			var rating float64
			var movieID int64
			var parseErr error

			movieID, parseErr = strconv.ParseInt(rowRecord[1], 10, 64)
			if parseErr != nil {
				continue
			}

			rating, parseErr = strconv.ParseFloat(rowRecord[2], 64)
			if parseErr != nil {
				continue
			}

			if _, ok := ratingsByMovieID[uint(movieID)]; !ok {
				ratingsByMovieID[uint(movieID)] = []float64{}
			}

			ratingsByMovieID[uint(movieID)] = append(ratingsByMovieID[uint(movieID)], rating)
		}

		return ratingsByMovieID, nil
	}
}

func LoadMoviesCSVFile(filepath string) (map[uint]*model.Movie, error) {
	if csvFile, err := os.Open(filepath); err == nil {
		reader := csv.NewReader(bufio.NewReader(csvFile))

		yearPattern, _ := regexp.Compile("\\(\\d{4}\\)")
		numericPattern, _ := regexp.Compile("\\d{4}")

		movieById := make(map[uint]*model.Movie)
		for {
			if row, readerErr := reader.Read(); readerErr != nil {
				if readerErr == io.EOF {
					break
				} else {
					fmt.Printf("Unexpected reader error: %v", readerErr)
				}
			} else {
				var id, year int64
				var parseErr error

				// Parameter bitSize defines range of values. If the value corresponding to s cannot be represented by a
				// signed integer of the given size, err.Err = ErrRange.
				id, parseErr = strconv.ParseInt(row[0], 10, 64)
				if parseErr != nil {
					continue
				}

				yearStr := yearPattern.FindString(row[1])
				trimmedTitle := strings.Trim(row[1], fmt.Sprintf(" %s", yearStr))

				year, parseErr = strconv.ParseInt(numericPattern.FindString(yearStr), 10, 64)
				if parseErr != nil {
					continue
				}

				movieById[uint(id)] = &model.Movie{
					ID:      uint(id),
					Year:    uint(year),
					Title:   trimmedTitle,
					Feature: pq.Float64Array{},
				}
			}
		}

		return movieById, nil
	} else {
		return nil, err
	}
}
