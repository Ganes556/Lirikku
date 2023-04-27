package models

import "gorm.io/gorm"

type SongLyric struct {
	Base
	UserID  uint      `json:"-" gorm:"type:unsignedInteger"`
	Artists []*Artist `json:"artists" gorm:"many2many:song_lyric_artists;constraint:OnDelete:CASCADE;"`
	Title   string    `json:"title" gorm:"type:varchar(150)"`
	Lyric   string    `json:"lyric"`
}

// hooks
func (s *SongLyric) AfterDelete(tx *gorm.DB) (err error) {
	tx.Model(&Artist{}).Unscoped().Delete(s.Artists)
	tx.Exec("ALTER TABLE artists AUTO_INCREMENT = 1")
	tx.Exec("ALTER TABLE song_lyrics AUTO_INCREMENT = 1")
	return
}

type SaveSongLyric struct {
	Title   string   `json:"title"`
	Lyric   string   `json:"lyric"`
	Artists []string `json:"artists"`
}

func (s *SaveSongLyric) Convert2SongLyric(userID uint) *SongLyric {
	var songLyric SongLyric
	songLyric.UserID = userID
	songLyric.Title = s.Title
	songLyric.Lyric = s.Lyric

	for _, artist := range s.Artists {
		var NewArtist Artist
		NewArtist.Name = artist
		songLyric.Artists = append(songLyric.Artists, &NewArtist)
	}

	return &songLyric
}