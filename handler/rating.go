// Copyright (c) 2018 Popcorn
// Author(s) Calvin Feng

package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"popcorn/model"
)

func NewRatingListHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var ratings []*model.Rating
		if err := db.Where("user_id = ?", vars["id"]).Find(&ratings).Error; err != nil {
			RenderError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if bytes, err := json.Marshal(ratings); err != nil {
			RenderError(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(bytes)
		}
	}
}

func NewRatingCreateHandler(db *gorm.DB, updateUserPreferenceQueue chan *model.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var rating model.Rating
		if err := decoder.Decode(&rating); err != nil {
			RenderError(w, "failed to parse request JSON into struct", http.StatusInternalServerError)
			return
		}

		if err := db.Create(&rating).Error; err != nil {
			RenderError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var user model.User
		if err := db.Where("id = ?", rating.UserID).Preload("Ratings").First(&user).Error; err == nil {
			// If the recommendation engine is too busy, we should just drop the request.
			select {
			case updateUserPreferenceQueue <- &user:
			default:
			}
		}

		if bytes, err := json.Marshal(&rating); err != nil {
			RenderError(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(bytes)
		}
	}
}
