// Copyright (c) 2018 Popcorn
// Author(s) Calvin Feng

package main

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"popcorn/handler"
	"popcorn/model"
)

func LoadRoutes(db *gorm.DB, updateUserPreferenceQueue chan *model.User) http.Handler {
	// Defining middleware
	logMiddleware := NewServerLoggingMiddleware()

	// Instantiate our router object
	muxRouter := mux.NewRouter().StrictSlash(true)

	// Name-spacing API
	api := muxRouter.PathPrefix("/api").Subrouter()

	// Sessions related
	api.Handle("/users/login", handler.NewSessionCreateHandler(db)).Methods("POST")
	api.Handle("/users/logout", handler.NewSessionDestroyHandler(db)).Methods("DELETE")
	api.Handle("/users/authenticate", handler.NewTokenAuthenticateHandler(db)).Methods("GET")

	// Users related
	api.Handle("/users/{id}/recommend", handler.NewPersonalizedRecommendationHandler(db)).Methods("GET")
	api.Handle("/users/register", handler.NewUserCreateHandler(db)).Methods("POST")
	api.Handle("/users/{id}/ratings", handler.NewRatingListHandler(db)).Methods("GET")
	api.Handle("/users", handler.NewUserListHandler(db)).Methods("GET")

	// Movies related
	api.Handle("/movies/popular", handler.NewPopularMovieListHandler(db)).Methods("GET")
	api.Handle("/movies/recommend", handler.NewMovieRecommendationHandler(db)).Methods("POST")
	api.Handle("/movies/details/{IMDBID}", handler.NewMovieDetailHandler(db)).Methods("GET")
	api.Handle("/movies/trailers/{IMDBID}", handler.NewMovieTrailerHandler(db)).Methods("GET")
	api.Handle("/movies", handler.NewMovieListHandler(db)).Methods("GET")

	// Ratings related
	api.Handle("/ratings", handler.NewRatingCreateHandler(db, updateUserPreferenceQueue)).Methods("POST")

	// Serve public folder to clients
	muxRouter.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	return logMiddleware(muxRouter)
}
