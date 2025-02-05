package models

type OvhSearchLyricResponse struct {
	Lyrics string `json:"lyrics"`
}

// func (o *OvhSearchLyricResponse) GetPubRes(shazamTrack ShazamTrack) PublicSongLyricResponse {
// 	return PublicSongLyricResponse{ArtistNames: shazamTrack.Subtitle, Title: shazamTrack.Title, Lyric: o.Lyrics}
// }

type OvhSearchTermResponse struct {
	Data  []DataSearchTerm `json:"data"`
	Total int              `json:"total"`
	Next  string           `json:"next"`
}

type DataSearchTerm struct {
	ID                    int    `json:"id"`
	Readable              bool   `json:"readable"`
	Title                 string `json:"title"`
	TitleShort            string `json:"title_short"`
	TitleVersion          string `json:"title_version"`
	Link                  string `json:"link"`
	Duration              int    `json:"duration"`
	Rank                  int    `json:"rank"`
	ExplicitLyrics        bool   `json:"explicit_lyrics"`
	ExplicitContentLyrics int    `json:"explicit_content_lyrics"`
	ExplicitContentCover  int    `json:"explicit_content_cover"`
	Preview               string `json:"preview"`
	Md5Image              string `json:"md5_image"`
	Artist                Artist `json:"artist"`
	Album                 Album  `json:"album"`
	Type                  string `json:"type"`
}

type Artist struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Link          string `json:"link"`
	Picture       string `json:"picture"`
	PictureSmall  string `json:"picture_small"`
	PictureMedium string `json:"picture_medium"`
	PictureBig    string `json:"picture_big"`
	PictureXl     string `json:"picture_xl"`
	Tracklist     string `json:"tracklist"`
	Type          string `json:"type"`
}

type Album struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Cover       string `json:"cover"`
	CoverSmall  string `json:"cover_small"`
	CoverMedium string `json:"cover_medium"`
	CoverBig    string `json:"cover_big"`
	CoverXl     string `json:"cover_xl"`
	Md5Image    string `json:"md5_image"`
	Tracklist   string `json:"tracklist"`
	Type        string `json:"type"`
}
