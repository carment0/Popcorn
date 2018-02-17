// Copyright (c) 2018 Popcorn
// Author(s) Calvin Feng

package model

import (
	"encoding/json"
	"github.com/lib/pq"
	"time"
)

// Since all matrix operations in gonum are done on float64, Postgres should also return a double precision float,
// despite that 64 bit float is an overkill for a rating that only has one decimal precision.
type Movie struct {
	// Model base class attributes
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	// Movie base attributes
	Title      string  `gorm:"type:varchar(500);index"          json:"title"`
	Year       uint    `gorm:"type:integer;index"               json:"year"`
	IMDBID     string  `gorm:"type:varchar(100);column:imdb_id" json:"imdb_id"`
	TMDBID     string  `gorm:"type:varchar(100);column:tmdb_id" json:"tmdb_id"`
	IMDBRating float64 `gorm:"type:float8;column:imdb_rating"   json:"-"`

	// NumRating is the number of ratings of this movie received from MovieLens dataset, while average rating is the
	// average of all the ratings received from MovieLens users.
	NumRating     int             `gorm:"type:integer"  json:"num_rating"`
	AverageRating float64         `gorm:"type:float8"   json:"average_rating"`
	Feature       pq.Float64Array `gorm:"type:float8[]" json:"-"`

	// The ratings here are submitted by the users of our web application, which is different from the ratings that came
	// from the MovieLens data set.
	Ratings []Rating `json:"-"`
}

type MovieDetail struct {
	IMDBID string          `gorm:"primary_key" json:"-"`
	Data   json.RawMessage `gorm:"type:bytea"  json:"-"`
}

type MovieTrailer struct {
	IMDBID string          `gorm:"primary_key" json:"-"`
	Data   json.RawMessage `gorm:"type:bytea"  json:"-"`
}
