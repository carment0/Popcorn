// Copyright (c) 2018 Popcorn
// Author(s) Calvin Feng

package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"gonum.org/v1/gonum/mat"
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

		for key, val := range reqData.Ratings {
			fmt.Printf("User rated movie %d with %.2f\n", key, val)
		}
		// Wait for Carmen's implementation...
		// The following logic does not actually recommend anything. It only finds the movie in the database and return
		// back whatever movie the user has rated. Also get rid of this as soon as possible. Making multiple SQL queries
		// is very inefficient.
		var movies []*model.Movie
		if err := db.Limit(10).Order("year desc").Find(&movies).Error; err != nil {
			RenderError(w, err.Error(), http.StatusInternalServerError)
			return
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
		vars := mux.Vars(r)

		var currentUser model.User
		if err := db.Where("id = ?", vars["id"]).First(&currentUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				RenderError(w, "user does not exist", http.StatusBadRequest)
				return
			}
			RenderError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var movies []*model.Movie
		if err := db.Order("id asc").Find(&movies).Error; err != nil {
			RenderError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		features := []float64{}
		for _, movie := range movies {
			features = append(features, movie.Feature...)
		}

		// K represents the feature dimension
		K := len(currentUser.Preference)

		U := mat.NewDense(1, K, currentUser.Preference)
		M := mat.NewDense(len(movies), K, features)

		predictedRatings := mat.NewDense(1, len(movies), nil)
		predictedRatings.Mul(U, M.T())

		// Fetch 10 recommendations
		recommendations := make([]*model.Movie, 0, 10)
		for j := 0; j < len(movies); j += 1 {
			if movies[j].Year < 2000 || movies[j].NumRating < 100 {
				continue
			}

			if predictedRatings.At(0, j) >= 4.0 {
				recommendations = append(recommendations, movies[j])
			}

			if len(recommendations) == 10 {
				break
			}
		}

		if bytes, err := json.Marshal(recommendations); err != nil {
			RenderError(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(bytes)
		}
	}
}
