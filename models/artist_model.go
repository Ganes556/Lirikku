package models

type Artist struct {
	Base
	SongLyrics []*SongLyric `gorm:"many2many:song_artists;" json:"song_artists"`
	Name       string       `json:"name" gorm:"type:varchar(150)"`
}
