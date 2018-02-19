// Copyright (c) 2018 Popcorn
// Author(s) Calvin Feng

// Package lowrank provides tools to perform low rank approximation on latent features of movies and users.
package lowrank

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/mat"
	"math/rand"
)

const TEST_RATIO = 0.10
const MAX_NUM_USER = 10000

type DataProcessor struct {
	MovieMap       map[int]*Movie
	UserIDToIndex  map[int]int
	UserIndexToID  map[int]int
	MovieIDToIndex map[int]int
	MovieIndexToID map[int]int

	// Training set is being fed into the algorithm
	TrainRatingMap map[int]map[int]float64

	// Test set is what we use to benchmark the accuracy of our model
	TestRatingMap map[int]map[int]float64
}

func NewDataProcessor(ratingFilePath string, movieFilepath string) (*DataProcessor, error) {
	rand.Seed(0)

	var movieMap map[int]*Movie
	var trainSet map[int]map[int]float64
	var testSet map[int]map[int]float64
	var loadErr error

	movieMap, loadErr = loadMovies(movieFilepath)
	if loadErr != nil {
		return nil, loadErr
	}

	trainSet, loadErr = loadRatingsByUserID(ratingFilePath, MAX_NUM_USER)
	if loadErr != nil {
		return nil, loadErr
	}

	var i, j int
	userIdToIndex := make(map[int]int)
	userIndexToId := make(map[int]int)
	movieIdToIndex := make(map[int]int)
	movieIndexToId := make(map[int]int)

	// Notice that we could have computed the average rating during loading the data from csv files. However, it would
	// have made the code less readable. I'd rather make the logic clear for this module. The performance gain does not
	// have a huge impact because it is only run once during server starts up.
	i = 0
	trainSetCount := 0
	testSetCount := 0

	testSet = make(map[int]map[int]float64)
	for userId := range trainSet {
		userIdToIndex[userId] = i
		userIndexToId[i] = userId
		i += 1

		if _, ok := testSet[userId]; !ok {
			testSet[userId] = make(map[int]float64)
		}

		for movieId := range trainSet[userId] {
			if rand.Float64() < TEST_RATIO {
				testSetCount += 1
				testSet[userId][movieId] = trainSet[userId][movieId]
			} else {
				trainSetCount += 1
				movieMap[movieId].Ratings = append(movieMap[movieId].Ratings, trainSet[userId][movieId])
			}
		}
	}

	// Compute average rating for each movie, while ignoring the data in test set
	j = 0
	for movieId := range movieMap {
		movieIdToIndex[movieId] = j
		movieIndexToId[j] = movieId
		j += 1

		movieMap[movieId].AvgRating = Average(movieMap[movieId].Ratings)
	}

	fmtString := "CSV data are loaded with %d training samples and %d test samples from %d users on %d movies"
	logMessage := fmt.Sprintf(fmtString, trainSetCount, testSetCount, len(userIdToIndex), len(movieIdToIndex))
	logrus.WithField("file", "lowrank.data_processor").Infof(logMessage)

	return &DataProcessor{
		MovieMap:       movieMap,
		UserIDToIndex:  userIdToIndex,
		UserIndexToID:  userIndexToId,
		MovieIDToIndex: movieIdToIndex,
		MovieIndexToID: movieIndexToId,
		TrainRatingMap: trainSet,
		TestRatingMap:  testSet,
	}, nil
}

// GetRatingMatrix returns a I by J matrix where I represents the ith user and J represents the jth movie. The rating
// matrix was supposed to be sparse but instead of filling it up with zero values. I've decided to set a movie's average
// rating as its baseline. All zero valued spaces will be filled by a movie's average rating.
func (dp *DataProcessor) GetRatingMatrix() *mat.Dense {
	I, J := len(dp.TrainRatingMap), len(dp.MovieMap)
	R := mat.NewDense(I, J, nil)
	for i := 0; i < I; i += 1 {
		for j := 0; j < J; j += 1 {
			userId := dp.UserIndexToID[i]
			movieId := dp.MovieIndexToID[j]

			isTrain := true
			if _, userExists := dp.TestRatingMap[userId]; userExists {
				if _, movieExists := dp.TestRatingMap[userId][movieId]; movieExists {
					isTrain = false
				}
			}

			if _, ok := dp.TrainRatingMap[userId][movieId]; ok && isTrain {
				R.Set(i, j, dp.TrainRatingMap[userId][movieId])
			} else {
				R.Set(i, j, dp.MovieMap[movieId].AvgRating)
			}
		}
	}

	return R
}
