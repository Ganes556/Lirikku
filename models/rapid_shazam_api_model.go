package models

// type RapidShazamSearchAudioResponse struct {
// 	Track struct {
// 		Key      string `json:"key"`
// 		Title    string `json:"title"`
// 		Subtitle string `json:"subtitle"`
// 		Sections []struct {
// 			Text []string `json:"text,omitempty"`
// 		} `json:"sections"`
// 	} `json:"track"`
// }

type RapidShazamSearchAudioResponse struct {
	Track     struct {
		Subtitle  string `json:"subtitle"`
		Title     string `json:"title"`
	} `json:"track"`
}

func (sa *RapidShazamSearchAudioResponse) GetInPublicSongLyricResponse() PublicSongsResponse {
	var res PublicSongsResponse
	res.Title = sa.Track.Title
	res.ArtistName = sa.Track.Subtitle
	return res
}