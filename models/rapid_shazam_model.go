package models

import "strings"

type ResponseRapidShazamSearchAudio struct {
	Track struct {
		Key      string `json:"key"`
		Title    string `json:"title"`
		Subtitle string `json:"subtitle"`
		Sections []struct {
			Text []string `json:"text,omitempty"`
		} `json:"sections"`
	} `json:"track"`
}

func (sa *ResponseRapidShazamSearchAudio) GetLyrics() string {
	var lyrics string
	if len(sa.Track.Sections) > 1 {
		lyrics = strings.Join(sa.Track.Sections[1].Text, "\n")
	}
	return lyrics
}

func (sa *ResponseRapidShazamSearchAudio) GetInResponsePublicSongLyric() ResponsePublicSongLyric {
	var res ResponsePublicSongLyric
	res.Title = sa.Track.Title
	res.ArtistNames = sa.Track.Subtitle
	res.Lyric = sa.GetLyrics()
	return res
}