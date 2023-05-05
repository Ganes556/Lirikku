package services

import (
	"github.com/Lirikku/models"
	"github.com/stretchr/testify/mock"
)

type MockMySongLyricsRepo struct {
	mock.Mock
}

func (m *MockMySongLyricsRepo) GetSongLyrics(userID uint, offset int) ([]models.SongLyricResponse, error) {
	args := m.Called(userID, offset)
	return args.Get(0).([]models.SongLyricResponse), args.Error(1)
}

func (m *MockMySongLyricsRepo) GetSongLyric(idSongLyric int,userID uint) (models.SongLyricResponse, error){
	args := m.Called(idSongLyric, userID)
	return args.Get(0).(models.SongLyricResponse), args.Error(1)
}

func (m *MockMySongLyricsRepo) CheckSongLyric(userID uint, req models.SongLyricWrite) error {
	args := m.Called(userID, req)
	return args.Error(0)
}

func (m *MockMySongLyricsRepo) SaveSongLyric(userID uint, req models.SongLyricWrite) error {
	args := m.Called(userID,req)
	return args.Error(0)
}

func (m *MockMySongLyricsRepo) SearchSongLyrics(userID uint, title, lyric, artist_names string, offset int) ([]models.SongLyricResponse,error){
	args := m.Called(userID, title, lyric, artist_names, offset)
	return args.Get(0).([]models.SongLyricResponse), args.Error(1)
}

func (m *MockMySongLyricsRepo) DeleteSongLyric(idSongLyric int, userID uint) error {
	args := m.Called(idSongLyric, userID)
	return args.Error(0)
}

func (m *MockMySongLyricsRepo) UpdateSongLyric(idSongLyric int, userID uint, req models.SongLyricWrite) error {
	args := m.Called(idSongLyric, userID, req)
	return args.Error(0)
}

