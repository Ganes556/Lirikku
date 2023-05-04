package services

import (
	"strconv"
	"sync"

	"github.com/Lirikku/models"
	"github.com/Lirikku/utils"
)

type IPublicSongLyricsService interface {
	SearchByTerm(term, types, limit string, offset int) ([]models.PublicSongLyricResponse, error)
	SearchByAudio(rawBases64 string) (models.PublicSongLyricResponse, error)
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

func (pub *PublicSongLyricsRepo) SearchByTerm(term, types, limit string, offset int) ([]models.PublicSongLyricResponse, error) {
	
	res, err := utils.RequestShazamSearchTerm(term, strconv.Itoa(offset), types, limit)

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

func (pub *PublicSongLyricsRepo) SearchByAudio(rawBases64 string) (models.PublicSongLyricResponse, error) {

	res, err := utils.RequestShazamSearchAudio(rawBases64)

	if err != nil {
		return models.PublicSongLyricResponse{}, err
	}

	if res.Track.Key == "" {
		return models.PublicSongLyricResponse{}, err
	}

	return res.GetInPublicSongLyricResponse(), nil

}