package models

import (
	"gorm.io/gorm"
)

type SongLyric struct {
	gorm.Model `json:"-"`
	UserID 	 uint `json:"user_id" gorm:"type:unsignedInteger"`
	Artists []*Artist `json:"artists" gorm:"many2many:song_artists;"`
	Title  string `json:"title" gorm:"type:varchar(150)"`
	Lyric  string `json:"lyric"`
}

type ReqSongLyric struct {
	Title  string `json:"title"`
	Lyric  string `json:"lyric"`
	Artists []string `json:"artists"`
}