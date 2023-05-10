package services

import (
	"errors"
	"sync"

	"github.com/Lirikku/models"
	"github.com/Lirikku/utils"
)

type IPublicSongLyricsService interface {
	SearchSongLyricsByTermShazam(term, types string, offset, pageSize int) ([]models.PublicSongLyricResponse, error)
	SearchSongLyricByAudioRapidShazam(rawBases64 string) (models.PublicSongLyricResponse, error)
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

func (pub *PublicSongLyricsRepo) SearchSongLyricsByTermShazam(term, types string, offset, pageSize int) ([]models.PublicSongLyricResponse, error) {

	res, err := utils.RequestShazamSearchTerm(term, types, offset, pageSize)

	if err != nil {
		return []models.PublicSongLyricResponse{}, err
	}

	keys := res.GetKeys()
	
	var resPublicSongLyrics = make([]models.PublicSongLyricResponse, len(keys))

	var wg sync.WaitGroup

	for i, key := range keys {
		wg.Add(1)
		go func (i int, key string) {
			defer wg.Done()

			res, err := utils.RequestShazamSearchKey(key)

			if err != nil {
				return
			}

			resPublicSongLyrics[i] = res.GetInPublicSongLyricResponse()
			
		}(i, key)

	}
	wg.Wait()

	return resPublicSongLyrics, nil
	
}

func (pub *PublicSongLyricsRepo) SearchSongLyricByAudioRapidShazam(rawBases64 string) (models.PublicSongLyricResponse, error) {

	res, err := utils.RequestShazamSearchAudio(rawBases64)
	
	if err != nil || res.Track.Key == "" {
		return models.PublicSongLyricResponse{}, errors.New("song lyric not found")
	}

	return res.GetInPublicSongLyricResponse(), nil

}