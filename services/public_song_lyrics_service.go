package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Lirikku/models"
	"github.com/Lirikku/utils"
)

type IPublicSongLyricsService interface {
	SearchSongsByTermShazam(term, types string, offset, pageSize int) ([]models.PublicSongsResponse, error)
	SearchSongLyricByAudioRapidShazam(rawBases64 string, q *models.ReqSearchAudio) (models.RapidShazamSearchAudioResponse, error)
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
			ArtistName: d.Track.Subtitle,
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
	artist = utils.ConvertUrl2Normal(artist)
	title = utils.ConvertUrl2Normal(title)
	
	delchunk := fmt.Sprintf("Paroles de la chanson %s par %s", title, artist)
	res.Lyrics = strings.ReplaceAll(res.Lyrics, delchunk, "")
	
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

func (pub *PublicSongLyricsRepo) SearchSongLyricByAudioRapidShazam(rawBases64 string, q *models.ReqSearchAudio) (models.RapidShazamSearchAudioResponse, error) {

	res, err := utils.RequestShazamSearchAudio(rawBases64, q)
	
	if err != nil {
		return res, errors.New("internal server error")
	}
	
	// if res.Track.Title == "" {
	// 	return models.PublicSongsResponse{}, errors.New("song lyric not found")
	// }

	return res, nil

	// return models.PublicSongsResponse{
	// 	ArtistName: res.Track.Subtitle,
	// 	Title: res.Track.Title,
	// }, nil
}
