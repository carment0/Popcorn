// Copyright (c) 2018 Popcorn
// Author(s) Calvin Feng, Carmen To

package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"gonum.org/v1/gonum/mat"
	"net/http"
	"popcorn/model"
	"strconv"
	"time"
	"math/rand"
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
		var maxYear uint = 2018
		var minYear uint = 1930

		if qMax := r.URL.Query().Get("max"); qMax != "" {
			if intValue, parseErr := strconv.ParseInt(qMax, 10, 64); parseErr == nil {
				maxYear = uint(intValue)
			}
		}

		if qMin := r.URL.Query().Get("min"); qMin != "" {
			if intValue, parseErr := strconv.ParseInt(qMin, 10, 64); parseErr == nil {
				minYear = uint(intValue)
			}
		}

		vars := mux.Vars(r)

		var currentUser model.User
		if err := db.Where("id = ?", vars["id"]).Preload("Ratings").First(&currentUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				RenderError(w, "user does not exist", http.StatusBadRequest)
				return
			}
			RenderError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rated := map[uint]bool{}
		for _, rating := range currentUser.Ratings {
			rated[rating.MovieID] = true
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

		rand.Seed(time.Now().UTC().UnixNano())
		// Fetch 10 recommendations
		recommendations := make([]*model.Movie, 0, 10)

		for len(recommendations) != 10 {
			j := rand.Intn(len(movies) - 1)
			if rated[movies[j].ID] {
				continue
			}

			if predictedRatings.At(0, j) < 3.0 {
				continue
			}

			if movies[j].Year >= minYear && movies[j].Year <= maxYear && movies[j].NumRating >= 20{
				recommendations = append(recommendations, movies[j])
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
