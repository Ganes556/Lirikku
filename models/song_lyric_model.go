package models

import (
	"gorm.io/gorm"
)

type SongLyric struct {
	Base
	UserID      uint   `json:",omitempty" gorm:"type:unsignedInteger"`
	ArtistNames string `json:"artist_names" gorm:"type:varchar(255)"`
	Title       string `json:"title" gorm:"type:varchar(150)"`
	Lyric       string `json:"lyric"`
}

type SongLyricWrite struct {
	ArtistNames string `json:"artist_names" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Lyric       string `json:"lyric" validate:"required"`
}

type SongLyricResponse struct {
	ID          uint   `json:"id"`
	ArtistNames string `json:"artist_names"`
	Title       string `json:"title"`
	Lyric       string `json:"lyric"`
}

type PublicSongsResponse struct {
	ArtistName string `json:"artist_names"`
	Title      string `json:"title"`
}

type PublicSongDetailResponse struct {
	ArtistName string `json:"artist_names"`
	Title      string `json:"title"`
	Lyric string `json:"lyric"`
}

// hooks
func (s *SongLyric) AfterDelete(tx *gorm.DB) (err error) {
	tx.Exec("ALTER TABLE song_lyrics AUTO_INCREMENT = 1")
	return
}
