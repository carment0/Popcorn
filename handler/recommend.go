// Copyright (c) 2018 Popcorn
// Author(s) Calvin Feng

package handler

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
	"popcorn/model"
)

type RecommendRequest struct {
	Ratings map[uint]float64
}

func NewMovieRecommendationHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var reqData RecommendRequest
		if err := decoder.Decode(&reqData); err != nil {
			RenderError(w, "Failed to parse request JSON into struct", http.StatusInternalServerError)
			return
		}

		// Wait for Carmen's implementation...
		// The following logic does not actually recommend anything. It only finds the movie in the database and return
		// back whatever movie the user has rated. Also get rid of this as soon as possible. Making multiple SQL queries
		// is very inefficient.
		movies := []*model.Movie{}
		for key, val := range reqData.Ratings {
			fmt.Printf("User rated movie %d with %.2f\n", key, val)
			movie := &model.Movie{}
			if err := db.Where("id = ?", key).First(movie).Error; err == nil {
				movies = append(movies, movie)
			}
		}

		if bytes, err := json.Marshal(movies); err != nil {
			RenderError(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(bytes)
		}
	}
}

func NewPersonalizedRecommendationHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implementation
	}
}
