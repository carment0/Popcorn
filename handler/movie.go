// Copyright (c) 2018 Popcorn
// Author(s) Carmen To, Calvin Feng

package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"io/ioutil"
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

func NewMovieRetrieveHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var movie model.Movie
		if err := db.Where("id = ?", vars["id"]).Find(&movie).Error; err != nil {
			RenderError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if bytes, err := json.Marshal(&movie); err != nil {
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
		if err := db.Limit(500).
			Order("num_rating desc").
			Order("average_rating desc").
			Find(&movies).Error; err != nil {
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

func NewMovieDetailHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var detail model.MovieDetail
		if err := db.Where("imdb_id = ?", vars["IMDBID"]).First(&detail).Error; err != nil {
			// If record is not found, make a request to The Movie Database and retrieve data.
			if err != gorm.ErrRecordNotFound {
				RenderError(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var httpError error
			var req *http.Request
			var res *http.Response

			req, httpError = http.NewRequest("GET", "https://api.themoviedb.org/3/movie/"+vars["IMDBID"], nil)
			if httpError != nil {
				RenderError(w, httpError.Error(), http.StatusInternalServerError)
				return
			}

			q := req.URL.Query()
			q.Add("api_key", "2afddf218bfb5d06ef460cc103af69bc")
			req.URL.RawQuery = q.Encode()

			res, httpError = http.DefaultClient.Do(req)
			if httpError != nil {
				RenderError(w, httpError.Error(), res.StatusCode)
				return
			}

			if res.StatusCode == http.StatusNotFound {
				RenderError(w, "there is no movie associated with the provided IMDB ID", res.StatusCode)
				return
			}

			bodyBytes, _ := ioutil.ReadAll(res.Body)
			TMDBResponse := struct {
				Status        int    `json:"status_code"`
				StatusMessage string `json:"status_message"`
			}{}
			// There are times when too many requests are made to TMDB. We need a catch statement here to make sure that
			// we don't store the error message into our database and treat it as actual movie detail data.
			if err := json.Unmarshal(bodyBytes, &TMDBResponse); err == nil {
				if TMDBResponse.Status == 25 {
					logrus.WithField("src", "handler.movie").Error("reaching request limits")
					RenderError(w, "request count is over the limit", http.StatusTooManyRequests)
					return
				}
			}

			detail = model.MovieDetail{
				IMDBID: vars["IMDBID"],
				Data:   bodyBytes,
			}

			if dbCreateErr := db.Create(&detail).Error; dbCreateErr != nil {
				RenderError(w, dbCreateErr.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write(detail.Data)
	}
}

func NewMovieTrailerHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var trailer model.MovieTrailer
		if err := db.Where("imdb_id = ?", vars["IMDBID"]).First(&trailer).Error; err != nil {
			// If record is not found, make a request to The Movie Database and retrieve data.
			if err != gorm.ErrRecordNotFound {
				RenderError(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var httpError error
			var req *http.Request
			var res *http.Response

			req, httpError = http.NewRequest(
				"GET",
				"https://api.themoviedb.org/3/movie/"+vars["IMDBID"]+"/videos",
				nil,
			)

			if httpError != nil {
				RenderError(w, httpError.Error(), http.StatusInternalServerError)
				return
			}

			q := req.URL.Query()
			q.Add("api_key", "2afddf218bfb5d06ef460cc103af69bc")
			req.URL.RawQuery = q.Encode()

			res, httpError = http.DefaultClient.Do(req)
			if httpError != nil {
				RenderError(w, httpError.Error(), res.StatusCode)
				return
			}

			if res.StatusCode == http.StatusNotFound {
				RenderError(w, "there is no movie associated with the provided IMDB ID", res.StatusCode)
				return
			}

			bodyBytes, _ := ioutil.ReadAll(res.Body)
			TMDBResponse := struct {
				Status        int    `json:"status_code"`
				StatusMessage string `json:"status_message"`
			}{}
			// There are times when too many requests are made to TMDB. We need a catch statement here to make sure that
			// we don't store the error message into our database and treat it as actual movie detail data.
			if err := json.Unmarshal(bodyBytes, &TMDBResponse); err == nil {
				if TMDBResponse.Status == 25 {
					logrus.WithField("src", "handler.movie").Error("reaching request limits")
					RenderError(w, "request count is over the limit", http.StatusTooManyRequests)
					return
				}
			}

			trailer = model.MovieTrailer{
				IMDBID: vars["IMDBID"],
				Data:   bodyBytes,
			}

			if dbCreateErr := db.Create(&trailer).Error; dbCreateErr != nil {
				RenderError(w, dbCreateErr.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write(trailer.Data)
	}
}
