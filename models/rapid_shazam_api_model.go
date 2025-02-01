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
	Track struct {
		Subtitle string `json:"subtitle"`
		Title    string `json:"title"`
	} `json:"track"`
	Tagid     string `json:"tagid"`
	Timestamp int64  `json:"timestamp"`
	Timezone  string `json:"timezone"`
}
