// Copyright (c) 2018 Popcorn
// Author(s) Carmen To

package handler

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

func NewTokenAuthenticateHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("session_token")
		if cookie == nil {
			RenderError(w, "cookie does not exist", http.StatusUnauthorized)
			return
		}

		if currentUser, err := FindUserByToken(db, cookie.Value); err == nil {
			if bytes, err := json.Marshal(currentUser); err != nil {
				RenderError(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write(bytes)
			}
		} else {
			RenderError(w, err.Error(), http.StatusUnauthorized)
		}
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewSessionCreateHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var reqData LoginRequest
		if err := decoder.Decode(&reqData); err != nil {
			RenderError(w, "Fail to parse request json into a struct", http.StatusInternalServerError)
			return
		}

		user, err := FindUserByCredential(db, reqData.Username, reqData.Password)
		if err != nil {
			RenderError(w, "Incorrect username/password combination", http.StatusUnauthorized)
			return
		}

		expiration := time.Now().Add(2 * 24 * time.Hour)
		cookie := http.Cookie{Name: "session_token", Value: user.SessionToken, Expires: expiration}
		http.SetCookie(w, &cookie)

		if bytes, err := json.Marshal(user); err != nil {
			RenderError(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(bytes)
		}

	}
}

type LogoutResponse struct {
	Username string `json:"username"`
}

func NewSessionDestroyHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("session_token")
		if currentUser, err := FindUserByToken(db, cookie.Value); err == nil {
			currentUser.ResetSessionToken()
			db.Save(currentUser)

			res := &LogoutResponse{currentUser.Username}

			if bytes, err := json.Marshal(res); err != nil {
				RenderError(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write(bytes)
			}
		} else {
			RenderError(w, "User is not found", http.StatusBadRequest)
		}
	}
}
