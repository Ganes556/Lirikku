package models

import "gorm.io/gorm"

type Artist struct {
	gorm.Model `json:"-"`
	SongLyrics      []*SongLyric `gorm:"many2many:song_artists;" json:"song_lyrics"`
	Name       string `json:"name" gorm:"type:varchar(150)"`
}
