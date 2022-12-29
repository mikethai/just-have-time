package model

import (
	"gorm.io/gorm"
)

type User struct {
	Msno      int64       `gorm:"primaryKey;unique"`
	Followers []User      `gorm:"many2many:follows;association_junction:Follow;foreignKey:msno"`
	Following []User      `gorm:"many2many:follows;association_junction:Follow;foreignKey:msno"`
	StorySong []StorySong `gorm:"foreignKey:msno;references:msno;constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

type Follow struct {
	gorm.Model       // Adds some metadata fields to the table
	FollowerID int64 `gorm:"index:follower_followee,unique;"`
	FolloweeID int64 `gorm:"index:follower_followee,unique;"`
}

type Song struct {
	SongID     string
	SongName   string
	ArtistName string
}

type Hashtag struct {
	gorm.Model             // Adds some metadata fields to the table
	Name       string      `json:"Value"`
	StorySong  []StorySong `gorm:"many2many:storysong_tag;foreignKey:id;"`
}

type StorySong struct {
	gorm.Model // Adds some metadata fields to the table
	SongID     string
	Msno       int64
	Song       Song      `gorm:"foreignKey:song_id;"`
	Hashtag    []Hashtag `gorm:"many2many:storysong_tag;foreignKey:id;"`
}
