// Copyright (c) 2018 Popcorn
// Author(s) Calvin Feng

package model

import "time"

type Rating struct {
	// Model base class attributes
	ID        uint      `gorm:"primary_key" json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	// Foreign Keys
	UserID  uint
	MovieID uint
	Value   float64 `gorm:"type:float8"`
}
