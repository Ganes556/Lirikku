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
	Artists []struct {
		Alias  string `json:"alias"`
	} `json:"artists"`
	Sections []struct {
		Text []string `json:"text,omitempty"`
	} `json:"sections"`
}

// search term

func (sm *ReponseShazamSearchTerm) GetKeys() (keys []string) {
	for _, track := range sm.Tracks.Hits {
		keys = append(keys, track.Track.Key)
	}
	return keys
}


// search key

func (sl *ResponseShazamSearchKey) GetLyrics() (lyrics string) {	
	if len(sl.Sections) > 1 {
		lyrics = strings.Join(sl.Sections[1].Text, "\n")
	}
	return lyrics
}

func (sl *ResponseShazamSearchKey) GetArtists() (artists string) {
	if len(sl.Artists) > 1 {
		for i, artist := range sl.Artists {
			if i == 0 {
				artists = artist.Alias
			} else {
				artists += ", " + artist.Alias
			}
		}
	}
	return artists
}

