// Copyright (c) 2018 Popcorn
// Author(s) Carmen To

package model

import (
	"time"
	"github.com/lib/pq"
)

type User struct {
	// Model base class attributes
	ID        uint      `gorm:"primary_key" json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	// User base attributes
	Username       string          `gorm:"type:varchar(100)"              json:"username"`
	SessionToken   string          `gorm:"type:varchar(100);unique_index" json:"session_token"`
	PasswordDigest []byte          `gorm:"type:bytea"                     json:"-"`
	Ratings        []Rating        `gorm:"ForeignKey:UserID"              json:"-"`
	Preference     pq.Float64Array `gorm:"type:float8[]"                  json:"-"`
}

func (u *User) ResetSessionToken() {
	if randStr, err := GenerateRandomString(20); err == nil {
		u.SessionToken = randStr
	}
}
