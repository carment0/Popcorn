// Copyright (c) 2018 Popcorn
// Author(s) Calvin Feng, Carmen To

package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"gonum.org/v1/gonum/mat"
	"math/rand"
	"net/http"
	"popcorn/model"
	"time"
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

type RecommendationRequestPayload struct {
	MaxYear    uint   `json:"max"`
	MinYear    uint   `json:"min"`
	Percentile uint   `json:"percent"`
	Skipped    []uint `json:"skipped"`
}

func NewPersonalizedRecommendationHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var payload RecommendationRequestPayload
		if err := decoder.Decode(&payload); err != nil {
			RenderError(w, "failed to parse request JSON into struct", http.StatusInternalServerError)
			return
		}
		fmt.Println(payload)

		var maxYear uint = 2018
		var minYear uint = 1930
		if payload.MaxYear != 0 {
			maxYear = payload.MaxYear
		}

		if payload.MinYear != 0 {
			minYear = payload.MinYear
		}

		// Find current user and get his/her ratings
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

		// Convert ratings into a set of rated movie ID's
		rated := map[uint]bool{}
		for _, rating := range currentUser.Ratings {
			rated[rating.MovieID] = true
		}

		// Convert skipped movies into a set of skipped movie ID's
		skipped := map[uint]bool{}
		for _, movieID := range payload.Skipped {
			skipped[movieID] = true
		}

		// Count how many movies there are in the database
		var movieCount int
		if err := db.Model(&model.Movie{}).Count(&movieCount).Error; err != nil {
			RenderError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var percentage float64
		switch payload.Percentile {
		case 100:
			// Give user the top 1% of most rated movies
			percentage = 0.01
		case 80:
			// Giver user the top 20% of most rated movies
			percentage = 0.2
		case 60:
			// Give user the top 40% of the most rated movies
			percentage = 0.4
		case 40:
			// Give user the top 60% of the most rated movies
			percentage = 0.6
		case 20:
			// Give user the top 80% of the most rated movies
			percentage = 0.8
		default:
			percentage = 1.0
		}

		limit := int(float64(movieCount) * percentage)

		var movies []*model.Movie
		if err := db.Limit(limit).
			Where("year >= ? and year <= ?", minYear, maxYear).
			Order("num_rating desc").
			Find(&movies).Error; err != nil {
			RenderError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Printf("fetched %d movies\n", len(movies))
		features := make([]float64, 0, len(movies))
		for _, movie := range movies {
			features = append(features, movie.Feature...)
		}

		// K represents the feature dimension
		K := len(currentUser.Preference)
		U := mat.NewDense(1, K, currentUser.Preference)
		M := mat.NewDense(len(movies), K, features)

		predictedRatings := mat.NewDense(1, len(movies), nil)
		predictedRatings.Mul(U, M.T())

		// Fetch 10 recommendations randomly
		rand.Seed(time.Now().UTC().UnixNano())
		recommendations := make([]*model.Movie, 0, 10)
		for len(recommendations) != 10 {
			j := rand.Intn(len(movies) - 1)
			if rated[movies[j].ID] {
				continue
			}

			if skipped[movies[j].ID] {
				continue
			}

			if predictedRatings.At(0, j) < 3.0 {
				continue
			}

			if movies[j].Year >= minYear && movies[j].Year <= maxYear && movies[j].NumRating >= 20 {
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
