package mocks

import (
	"github.com/Lirikku/models"
	"github.com/stretchr/testify/mock"
)

type MockPublicSongLyricsRepo struct {
	mock.Mock
}


func (m *MockPublicSongLyricsRepo)	SearchSongLyricsByTermShazam(term, types string, offset, pageSize int) ([]models.PublicSongDetailResponse, error){
	args := m.Called(term, types, offset, pageSize)
	return args.Get(0).([]models.PublicSongDetailResponse), args.Error(1)
	
}
func (m *MockPublicSongLyricsRepo)	SearchSongLyricByAudioRapidShazam(rawBases64 string) (models.PublicSongDetailResponse, error){
	args := m.Called(rawBases64)
	return args.Get(0).(models.PublicSongDetailResponse), args.Error(1)
}
