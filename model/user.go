package model

import "time"

type User struct {
	// Model base class attributes
	ID        uint      `gorm:"primary_key" json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	// User base attributes
	Username       string   `gorm:"type:varchar(100)"              json:"username"`
	SessionToken   string   `gorm:"type:varchar(100);unique_index" json:"-"`
	PasswordDigest []byte   `gorm:"type:bytea"                     json:"-"`
	Ratings        []Rating `gorm:"ForeignKey:UserID"              json:"ratings"`
}

func (u *User) ResetSessionToken() {
	if randStr, err := GenerateRandomString(20); err == nil {
		u.SessionToken = randStr
	}
}
