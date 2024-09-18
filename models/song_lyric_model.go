package models

import (
	"gorm.io/gorm"
)

type SongLyric struct {
	Base
	UserID      uint   `json:",omitempty" gorm:"type:unsignedInteger"`
	ArtistNames string `json:"artist_names" gorm:"type:varchar(255)"`
	KeyShazam   string `gorm:"type:varchar(50);uniqueIndex"`
	Title       string `json:"title" gorm:"type:varchar(150)"`
	Lyric       string `json:"lyric"`
}

type SongLyricWrite struct {
	ID          string `json:"id"`
	ArtistNames string `json:"artist_names"`
	Key         string `json:"key,omitempty"`
	Title       string `json:"title"`
	Lyric       string `json:"lyric"`
}

type SongLyricResponse struct {
	ID          uint   `json:"id"`
	ArtistNames string `json:"artist_names"`
	Title       string `json:"title"`
	Lyric       string `json:"lyric"`
}

type PublicSongsGetByAudioBase64 struct {
	AudioBase64 string `json:"audio_base64"`
}

type PublicSaveSong struct {
	ArtistNames string `json:"artist_names"`
	Key         string `json:"key"`
	Title       string `json:"title"`
}

type PublicSongsResponse struct {
	ArtistNames string `json:"artist_names"`
	Key         string `json:"key"`
	Saved       bool   `json:"saved"`
	Title       string `json:"title"`
}

type PublicSongDetailResponse struct {
	ArtistNames string `json:"artist_names"`
	Title       string `json:"title"`
	Lyric       string `json:"lyric"`
}

// hooks
func (s *SongLyric) AfterDelete(tx *gorm.DB) (err error) {
	tx.Exec("ALTER TABLE song_lyrics AUTO_INCREMENT = 1")
	return
}
