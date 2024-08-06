package models

import (
	"strings"
)

type ShazamTrack struct {
	Key      string `json:"key"`
	Subtitle string `json:"subtitle"`
	Title    string `json:"title"`
}

type ShazamSearchTermResponse struct {
	Tracks struct {
		Hits []struct {
			Track ShazamTrack `json:"track"`
		} `json:"hits"`
	} `json:"tracks"`
}

type ShazamSearchKeyResponse struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Sections []struct {
		Text []string `json:"text,omitempty"`
	} `json:"sections"`
}

// search term

func (sm *ShazamSearchTermResponse) GetKeys() []string {
	var keys []string
	for _, track := range sm.Tracks.Hits {
		keys = append(keys, track.Track.Key)
	}
	return keys
}

// search key

func (sk *ShazamSearchKeyResponse) GetLyric() string {
	var lyrics string
	if len(sk.Sections) > 1 {
		lyrics = strings.Join(sk.Sections[1].Text, "\n")
	}
	return lyrics
}

// func (sk *ShazamSearchKeyResponse) GetInPublicSongLyricResponse() PublicSongLyricResponse {
// 	var res PublicSongLyricResponse
// 	res.Title = sk.Title
// 	res.ArtistName = sk.Subtitle
// 	res.Lyric = sk.GetLyric()
// 	return res
// }
