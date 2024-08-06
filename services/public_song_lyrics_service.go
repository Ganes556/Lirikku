package services

import (
	"errors"

	"github.com/Lirikku/models"
	"github.com/Lirikku/utils"
)

type IPublicSongLyricsService interface {
	SearchSongsByTermShazam(term, types string, offset, pageSize int) ([]models.PublicSongsResponse, error)
	SearchSongLyricByAudioRapidShazam(rawBases64 string) (models.PublicSongDetailResponse, error)
	GetSongDetail(artist string, title string) (models.PublicSongDetailResponse, error)
	SearchTermByOvh(term string) (models.OvhSearchTermResponse, error)
	
}

type PublicSongLyricsRepo struct{}

var publicSongLyricsRepo IPublicSongLyricsService

func init() {
	publicSongLyricsRepo = &PublicSongLyricsRepo{}
}

func GetPublicSongLyricsRepo() IPublicSongLyricsService {
	return publicSongLyricsRepo
}

func SetPublicSongLyricsRepo(repo IPublicSongLyricsService) {
	publicSongLyricsRepo = repo
}

func (pub *PublicSongLyricsRepo) SearchSongsByTermShazam(term, types string, offset, pageSize int) ([]models.PublicSongsResponse, error) {

	res, err := utils.RequestTermByShazam(term, types, offset, pageSize)

	if err != nil {
		return nil, err
	}

	var resPub = make([]models.PublicSongsResponse, len(res.Tracks.Hits))
	for i, d := range res.Tracks.Hits {
		resPub[i] = models.PublicSongsResponse{
			ArtistName: utils.ConvertCapitalize(d.Track.Subtitle),
			Title: d.Track.Title,
		}
	}
	return resPub, nil
}

func (pub *PublicSongLyricsRepo) GetSongDetail(artist string, title string) (models.PublicSongDetailResponse, error) {
	res, err := utils.RequestLyricByOvh(artist, title)
	if err != nil {
		return models.PublicSongDetailResponse{}, err
	}
	return models.PublicSongDetailResponse{
		ArtistName: artist,
		Title: title,
		Lyric: res.Lyrics,
	}, nil
}

func (pub *PublicSongLyricsRepo) SearchTermByOvh(term string) (models.OvhSearchTermResponse, error) {
	res, err := utils.RequestTermByOvh(term)
	if err != nil {
		return models.OvhSearchTermResponse{}, err
	}
	return res, nil
}

func (pub *PublicSongLyricsRepo) SearchLyricByOvh(artist string, title string) (models.OvhSearchLyricResponse, error) {
	res, err := utils.RequestLyricByOvh(artist, title)
	if err != nil {
		return models.OvhSearchLyricResponse{}, err
	}
	return res, nil
}

func (pub *PublicSongLyricsRepo) SearchSongLyricByAudioRapidShazam(rawBases64 string) (models.PublicSongDetailResponse, error) {

	res, err := utils.RequestShazamSearchAudio(rawBases64)

	if err != nil || res.Track.Key == "" {
		return models.PublicSongDetailResponse{}, errors.New("song lyric not found")
	}

	return res.GetInPublicSongLyricResponse(), nil

}
