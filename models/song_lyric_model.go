package models

type SongLyric struct {
	Base
	UserID  uint      `json:"user_id" gorm:"type:unsignedInteger"`
	Artists []*Artist `json:"artists" gorm:"many2many:song_artists;"`
	Title   string    `json:"title" gorm:"type:varchar(150)"`
	Lyric   string    `json:"lyric"`
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