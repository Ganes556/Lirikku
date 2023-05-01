package models

type ResponseRapidShazamSearchAudio struct {
	Track struct {
		Key string `json:"key"`
	} `json:"track"`
}
