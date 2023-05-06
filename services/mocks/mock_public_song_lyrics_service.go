package mocks

import (
	"github.com/Lirikku/models"
	"github.com/stretchr/testify/mock"
)

type MockPublicSongLyricsRepo struct {
	mock.Mock
}


func (m *MockPublicSongLyricsRepo)	SearchSongLyricsByTermShazam(term, types, limit string, offset int) ([]models.PublicSongLyricResponse, error){
	args := m.Called(term, types, limit, offset)
	return args.Get(0).([]models.PublicSongLyricResponse), args.Error(1)
	
}
func (m *MockPublicSongLyricsRepo)	SearchSongLyricByAudioRapidShazam(rawBases64 string) (models.PublicSongLyricResponse, error){
	args := m.Called(rawBases64)
	return args.Get(0).(models.PublicSongLyricResponse), args.Error(1)
}
