// Copyright (c) 2018 Popcorn
// Author(s) Calvin Feng

// Package lowrank provides tools to perform low rank approximation on latent features of movies and users.
package lowrank

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Movie struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	AvgRating float64   `json:"avg_rating"`
	Ratings   []float64 `json:"ratings"`
}

// loadMovies will return a map that maps movie id to the title of the movie. We don't really care about meta
// information of the movies. We only want to know the mapping of matrix indices to movie ids.
func loadMovies(filepath string) (map[int]*Movie, error) {
	if csvFile, fileErr := os.Open(filepath); fileErr != nil {
		return nil, fileErr
	} else {
		reader := csv.NewReader(bufio.NewReader(csvFile))
		movieByID := make(map[int]*Movie)
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

			id, parseErr := strconv.ParseInt(rowRecord[0], 10, 64)
			if parseErr != nil {
				continue
			}

			movieByID[int(id)] = &Movie{
				ID:    int(id),
				Title: rowRecord[1],
			}
		}

		return movieByID, nil
	}
}

// loadRatingsByUserID will return a map that maps user id to a map of movie id to rating.
// Example:
// {
//   userID: {
//     movieID: 4.5,
//     movieID: 3.5
//   }
// }
func loadRatingsByUserID(filepath string, userMaxNum int) (map[int]map[int]float64, error) {
	if csvFile, fileErr := os.Open(filepath); fileErr != nil {
		return nil, fileErr
	} else {
		reader := csv.NewReader(bufio.NewReader(csvFile))
		ratingsByUserID := make(map[int]map[int]float64)
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

			var userID, movieID int64
			var rating float64
			var parseErr error

			userID, parseErr = strconv.ParseInt(rowRecord[0], 10, 64)
			if parseErr != nil {
				continue
			}

			movieID, parseErr = strconv.ParseInt(rowRecord[1], 10, 64)
			if parseErr != nil {
				continue
			}

			rating, parseErr = strconv.ParseFloat(rowRecord[2], 64)
			if parseErr != nil {
				continue
			}

			if _, ok := ratingsByUserID[int(userID)]; !ok {
				ratingsByUserID[int(userID)] = make(map[int]float64)
			}

			ratingsByUserID[int(userID)][int(movieID)] = rating
		}

		// We are only interested in users who have rated at least 300 movies.
		reducedMap := make(map[int]map[int]float64)
		for userId := range ratingsByUserID {
			if len(ratingsByUserID[userId]) > 300 {
				reducedMap[userId] = ratingsByUserID[userId]
			}

			if len(reducedMap) == userMaxNum {
				break
			}
		}

		return reducedMap, nil
	}
}
