package models

import (
	"strings"
)

type ReponseShazamSearchTerm struct {
	Tracks struct {
		Hits []struct {
			Track struct {
				Key   string `json:"key"`
			} `json:"track"`
		} `json:"hits"`
	} `json:"tracks"`
}

type ResponseShazamSearchKey struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Sections []struct {
		Text []string `json:"text,omitempty"`
	} `json:"sections"`
}

// search term

func (sm *ReponseShazamSearchTerm) GetKeys() []string {
	var keys []string
	for _, track := range sm.Tracks.Hits {
		keys = append(keys, track.Track.Key)
	}
	return keys
}


// search key

func (sk *ResponseShazamSearchKey) GetLyrics() string {	
	var lyrics string
	if len(sk.Sections) > 1 {
		lyrics = strings.Join(sk.Sections[1].Text, "\n")
	}
	return lyrics
}


func (sk *ResponseShazamSearchKey) GetInResponsePublicSongLyric() ResponsePublicSongLyric {
	var res ResponsePublicSongLyric
	res.Title = sk.Title
	res.ArtistNames = sk.Subtitle
	res.Lyric = sk.GetLyrics()
	return res
}