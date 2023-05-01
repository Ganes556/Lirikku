package models

import (
	"encoding/json"
	"io"
	"strings"
)

type ShazamMetadata struct {
	Tracks struct {
		Hits []struct {
			Track struct {
				Key   string `json:"key"`
				Title string `json:"title"`
			} `json:"track"`
		} `json:"hits"`
	} `json:"tracks"`
	Artists struct {
		Hits []struct {
			Artist struct {
				Name string `json:"name"`
			} `json:"artist"`
		} `json:"hits"`
	} `json:"artists"`
}

type ShazamLyric struct {
	Sections []struct {
		Text []string `json:"text,omitempty"`
	} `json:"sections"`
}

func (sm *ShazamMetadata) GetTitles() (titles []string) {
	if len(sm.Tracks.Hits) > 0 {
		for _, track := range sm.Tracks.Hits {
			titles = append(titles, track.Track.Title)
		}
	}
	return titles
}

func (sm *ShazamMetadata) GetArtists() (artists []string) {
	if len(sm.Artists.Hits) > 0 {
		for _, artist := range sm.Artists.Hits {
			artists = append(artists, artist.Artist.Name)
		}
	}
	return artists
}

func (sm *ShazamMetadata) GetKeys() (keys []string) {
	for _, track := range sm.Tracks.Hits {
		keys = append(keys, track.Track.Key)
	}
	return keys
}

func (sl *ShazamLyric) Decode(body io.Reader) {
	json.NewDecoder(body).Decode(&sl)
}

func (sl *ShazamLyric) GetLyrics() (lyrics string) {	
	if len(sl.Sections) > 1 {
		lyrics = strings.Join(sl.Sections[1].Text, "\n")
	}
	return lyrics
}