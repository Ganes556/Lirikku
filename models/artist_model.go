package models

import "gorm.io/gorm"

type Artist struct {
	gorm.Model `json:"-"`
	Songs      []SongLyric `gorm:"many2many:song_artists;" json:"songs"`
	Name       string `json:"name" gorm:"type:varchar(150)"`
}
