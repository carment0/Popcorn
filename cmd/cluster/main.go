// Copyright (c) 2018 Popcorn
// Author(s) Carmen To
package main

import (
	"fmt"
	"math/rand"
	"popcorn/kmeans"
	"time"
)

const DIR = "datasets/production/"

func main() {
	rand.Seed(time.Now().Unix())

	movies, err := kmeans.ReadFromCSV(DIR + "features.csv")
	if err != nil {
		fmt.Println("Failed to read CSV", err)
	} else {
		assignedMovies := kmeans.MovieClustering(movies)
		kmeans.WriteToCSV(DIR + "clusters.csv", assignedMovies)
	}
}
