package models

type Artist struct {
	Base
	SongLyrics []*SongLyric `json:"song_artists" gorm:"many2many:song_lyrics_artists;"`
	Name       string       `json:"name" gorm:"type:varchar(150)"`
}
