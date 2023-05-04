package services

import (
	"github.com/Lirikku/configs"
	"github.com/Lirikku/models"
)

type IMySongLyricsService interface {
	GetMySongLyrics(userID uint, offset int) ([]models.SongLyricResponse, error)
	GetMySongLyric(idSongLyric int,userID uint) (models.SongLyricResponse, error)
	SaveMySongLyric(userID uint, req models.SongLyricWrite) error
	SearchMySongLyrics(userID uint, title, lyric, artist_names string, offset int) ([]models.SongLyricResponse,error)
	DeleteMySongLyric(idSongLyric int, userID uint) error
	UpdateMySongLyric(idSongLyric int, userID uint, req models.SongLyricWrite) error
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

func (my *MySongLyricsRepo) GetMySongLyrics(userID uint, offset int) ([]models.SongLyricResponse, error) {
	var res []models.SongLyricResponse
	err:= configs.DB.Model(&models.SongLyric{}).Limit(5).Offset(offset).Find(&res, "user_id = ?", userID).Error
	if err != nil {
		return res, err
	}
	return res, nil
}

func (my *MySongLyricsRepo) GetMySongLyric(idSongLyric int,userID uint) (models.SongLyricResponse, error){
	var res models.SongLyricResponse
	
	err := configs.DB.Model(&models.SongLyric{}).First(&res, "id = ? AND user_id = ?", idSongLyric, userID).Error
	if err != nil {
		return res, err
	}
	return res, nil
}


func (my *MySongLyricsRepo) SaveMySongLyric(userID uint, req models.SongLyricWrite) error {
	newSongLyric := models.SongLyric{
		UserID: userID,
		Title: req.Title,
		Lyric: req.Lyric,
		ArtistNames: req.ArtistNames,
	}

	if err := configs.DB.Create(&newSongLyric).Error; err != nil {
		return err
	}

	return nil
}

func (my *MySongLyricsRepo) SearchMySongLyrics(userID uint, title, lyric, artist_names string, offset int) ([]models.SongLyricResponse,error) {

	var resSongLyrics []models.SongLyricResponse

	err := configs.DB.Model(&models.SongLyric{}).Where("user_id = ? AND title LIKE ? AND lyric LIKE ? AND artist_names LIKE ?", userID, "%"+title+"%", "%"+lyric+"%", "%"+artist_names+"%").Limit(5).Offset(offset).Find(&resSongLyrics).Error

	if err != nil {
		return resSongLyrics, err
	}

	return resSongLyrics, nil

}

func (my *MySongLyricsRepo) DeleteMySongLyric(idSongLyric int, userID uint) error{
	
	err := configs.DB.Unscoped().Delete(&models.SongLyric{}, "id = ? AND user_id = ?", idSongLyric, userID).Error

	if err != nil {
		return err
	}
	return  nil
}

func (my *MySongLyricsRepo) UpdateMySongLyric(idSongLyric int, userID uint, req models.SongLyricWrite) error{
	updateSongLyric := models.SongLyric{
		Title: req.Title,
		Lyric: req.Lyric,
	}
	err := configs.DB.Model(&models.SongLyric{}).Where("id = ? AND user_id = ?", idSongLyric, userID).Updates(updateSongLyric).Error

	if err != nil {
		return err
	}

	return  nil
}

