package model

import (
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type User struct {
	Msno         string `gorm:"primaryKey;unique"`
	MsnoInt      int64
	UserHashTags pgtype.JSONB `gorm:"type:jsonb;default:'[]';not null"`
	Followers    []User       `gorm:"many2many:follows;association_junction:Follow;foreignKey:msno"`
	Following    []User       `gorm:"many2many:follows;association_junction:Follow;foreignKey:msno"`
}

type Follow struct {
	gorm.Model        // Adds some metadata fields to the table
	FollowerID string `gorm:"index:follower_followee,unique;"`
	FolloweeID string `gorm:"index:follower_followee,unique;"`
}

type Song struct {
	SongID     string
	SongName   string
	ArtistName string
}

type Hashtag struct {
	gorm.Model             // Adds some metadata fields to the table
	Name       string      `gorm:"index:tag_name_count" json:"Value"`
	StorySong  []StorySong `gorm:"many2many:storysong_tag;foreignKey:id;"`
	Count      uint64      `gorm:"index:tag_name_count;default:1;not null"`
}

type StorySong struct {
	gorm.Model // Adds some metadata fields to the table
	SongID     string
	Msno       string
	Song       Song      `gorm:"foreignKey:song_id;"`
	Hashtag    []Hashtag `gorm:"many2many:storysong_tag;foreignKey:id;"`
	User       User      `gorm:"foreignKey:msno;references:msno;"`
}
