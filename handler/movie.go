// Copyright (c) 2018 Popcorn
// Author(s) Carmen To

package handler

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"net/http"
	"popcorn/model"
)

func NewMovieListHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var movies []*model.Movie
		if err := db.Order("year desc").Find(&movies).Error; err != nil {
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

func NewPopularMovieListHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var movies []*model.Movie
		if err := db.Limit(300).Order("num_rating desc").Order("average_rating desc").Find(&movies).Error; err != nil {
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
