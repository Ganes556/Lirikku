package models

import "gorm.io/gorm"

type SongLyric struct {
	ID      uint `json:"id" gorm:"primarykey"` 
	UserID  uint      `json:",omitempty" gorm:"type:unsignedInteger"`
	Artists []*Artist `json:"artists" gorm:"many2many:song_lyric_artists;constraint:OnDelete:CASCADE;"`
	Title   string    `json:"title" gorm:"type:varchar(150)"`
	Lyric   string    `json:"lyric"`
	Base
}

// save & update
type WriteSongLyric struct {
	Title   string   `json:"title"`
	Lyric   string   `json:"lyric"`
	Artists []string `json:"artists"`
}

type ResponseSongLyric struct {
	ID 		uint     `json:"id"`
	Artists []string `json:"artists"`
	Title   string   `json:"title"`
	Lyric   string   `json:"lyric"`
}

// hooks
func (s *SongLyric) AfterDelete(tx *gorm.DB) (err error) {
	tx.Model(&Artist{}).Unscoped().Delete(s.Artists)
	tx.Exec("ALTER TABLE artists AUTO_INCREMENT = 1")
	tx.Exec("ALTER TABLE song_lyrics AUTO_INCREMENT = 1")
	return
}

func (s *SongLyric) Convert2ResSongLyric() ResponseSongLyric {
	var resSongLyric ResponseSongLyric
	resSongLyric.ID = s.ID
	resSongLyric.Title = s.Title
	resSongLyric.Lyric = s.Lyric

	for _, artist := range s.Artists {
		resSongLyric.Artists = append(resSongLyric.Artists, artist.Name)
	}

	return resSongLyric
}


func (ws *WriteSongLyric) Convert2SongLyric(userID uint) SongLyric {
	var songLyric SongLyric
	songLyric.UserID = userID
	songLyric.Title = ws.Title
	songLyric.Lyric = ws.Lyric

	for _, artist := range ws.Artists {
		var NewArtist Artist
		NewArtist.Name = artist
		songLyric.Artists = append(songLyric.Artists, &NewArtist)
	}

	return songLyric
}