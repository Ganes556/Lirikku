package services

import (
	"errors"
	"slices"

	"github.com/Lirikku/configs"
	"github.com/Lirikku/models"
)

type IMySongLyricsService interface {
	GetSongLyrics(userID uint, offset, pageSize int) ([]*models.SongLyricResponse, error)
	GetSongLyric(idSongLyric int,userID uint) (models.SongLyricResponse, error)
	CheckSongLyric(userID uint, req models.SongLyricWrite) error 
	SaveSongLyric(userID uint, req models.SongLyricWrite) error
	SearchSongLyrics(userID uint, term string, offset, pageSize int) ([]*models.SongLyricResponse,error)
	DeleteSongLyric(idSongLyric int, userID uint) error
	UpdateSongLyric(idSongLyric int, userID uint, req models.SongLyricUpdate) error
	CheckKeyPub(userId uint, pubRes []*models.PublicSongsResponse, keyPub []string) error
}

type MySongLyricsRepo struct{}

var mySongLyricsRepo IMySongLyricsService

func init() {
	mySongLyricsRepo = &MySongLyricsRepo{}
}

func GetMySongLyricsRepo() IMySongLyricsService {
	return mySongLyricsRepo
}

func SetMySongLyricsRepo(repo IMySongLyricsService) {
	mySongLyricsRepo = repo
}

func (my *MySongLyricsRepo) GetSongLyrics(userID uint, offset, pageSize int) ([]*models.SongLyricResponse, error) {
	var res []*models.SongLyricResponse
	
	err := configs.DB.Model(&models.SongLyric{}).Omit("lyric").Limit(pageSize).Offset(offset).Find(&res, "user_id = ?", userID).Error
	
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (my *MySongLyricsRepo) GetSongLyric(idSongLyric int, userID uint) (models.SongLyricResponse, error){
	var res models.SongLyricResponse
	
	err := configs.DB.Model(&models.SongLyric{}).First(&res, "id = ? AND user_id = ?", idSongLyric, userID).Error
	if err != nil {
		return res, err
	}
	return res, nil
}

func (my *MySongLyricsRepo) CheckSongLyric(userID uint, req models.SongLyricWrite) error {
	
	err := configs.DB.First(&models.SongLyric{}, "user_id = ? AND title = ? AND lyric = ? AND artist_names = ?", userID, req.Title, req.Lyric, req.ArtistNames).Error

	if err == nil {
		return errors.New("song lyric already saved")
	}

	return nil
	
}

func (my *MySongLyricsRepo) SaveSongLyric(userID uint, req models.SongLyricWrite) error {
	newSongLyric := models.SongLyric{
		UserID: userID,
		Title: req.Title,
		Lyric: req.Lyric,
		KeyShazam: req.Key,
		ArtistNames: req.ArtistNames,
	}

	if err := configs.DB.Create(&newSongLyric).Error; err != nil {
		return err
	}

	return nil
}

func (my *MySongLyricsRepo) SearchSongLyrics(userID uint, term string, offset, pageSize int) ([]*models.SongLyricResponse,error) {

	var resSongLyrics []*models.SongLyricResponse

	err := configs.DB.Model(&models.SongLyric{}).Where("user_id = ? AND (title LIKE ? OR lyric LIKE ? OR artist_names LIKE ?)", userID, "%"+term+"%", "%"+term+"%", "%"+term+"%").Limit(pageSize).Offset(offset).Find(&resSongLyrics).Error

	if err != nil {
		return nil, err
	}

	return resSongLyrics, nil

}

func (my *MySongLyricsRepo) CheckKeyPub(userId uint, pubRes []*models.PublicSongsResponse, keyPub []string) error {
	var savedKeySong = make([]string, len(keyPub))
	err := configs.DB.Select("key_shazam").Model(&models.SongLyric{}).Where("user_id = ? AND key_shazam IN (?)", userId, keyPub).Pluck("key", &savedKeySong).Error
	if err != nil {
		return err
	}
	
	for i, v := range keyPub {
		if slices.Contains[[]string](savedKeySong, v) {
			pubRes[i].Saved = true
		}
	}
	return nil
}

func (my *MySongLyricsRepo) DeleteSongLyric(idSongLyric int, userID uint) error{
	
	err := configs.DB.Unscoped().Delete(&models.SongLyric{}, "id = ? AND user_id = ?", idSongLyric, userID).Error

	if err != nil {
		return err
	}
	return  nil
}

func (my *MySongLyricsRepo) UpdateSongLyric(idSongLyric int, userID uint, req models.SongLyricUpdate) error{
	
	updateSongLyric := models.SongLyric{
		Title: req.Title,
		ArtistNames: req.ArtistNames,
		Lyric: req.Lyric,
	}

	err := configs.DB.Model(&models.SongLyric{}).Where("id = ? AND user_id = ?", idSongLyric, userID).Updates(updateSongLyric).Error

	if err != nil {
		return err
	}

	return  nil
}

