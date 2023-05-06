package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/Lirikku/models"
	"github.com/Lirikku/services"
	mockService "github.com/Lirikku/services/mocks"
	"github.com/Lirikku/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSearchTermSongLyrics(t *testing.T) {
	mockPubSongLyricsRepo := mockService.MockPublicSongLyricsRepo{}	
	services.SetPublicSongLyricsRepo(&mockPubSongLyricsRepo)
	pubSongLyricsController := NewPublicSongLyricsController(&mockPubSongLyricsRepo)
	
	tests := []struct {
		name         string
		param 			 string
		expectedBody echo.Map
		expectedCode int
		wantErr      bool
	}{
		{
			name: "Success",
			param: "?term=test&offset=0",
			expectedCode: http.StatusOK,
			expectedBody: echo.Map{
				"next": "",
				"public_song_lyrics": []models.PublicSongLyricResponse{
					{
						Title: "test",
						ArtistNames: "test",
						Lyric: "test",
					},
				},
			},
			wantErr:      false,
		},
		{
			name: "Failed: offset must be a number",
			param: "?term=test&offset=test",
			expectedCode: http.StatusBadRequest,
			expectedBody: echo.Map{
				"message": "offset must be a number",
			},
			wantErr:      true,
		},
		{
			name: "Failed: internal server error (SearchSongLyricsByTermShazam)",
			param: "?term=test&offset=0",
			expectedCode: http.StatusInternalServerError,
			expectedBody: echo.Map{
				"message": "internal server error",
			},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/song_lyrics/public/search" + tt.param, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			
			term := c.QueryParam("term")
			offset := c.QueryParam("offset")

			offsetInt, err := utils.CheckOffset(offset)

			if tt.name != "Failed: offset must be a number" {
				assert.NoError(t, err)
			}
			
			var data []models.PublicSongLyricResponse

			if tt.wantErr {
				mockPubSongLyricsRepo.On("SearchSongLyricsByTermShazam", term, "artists,songs", "5", offsetInt).Return(data, errors.New(tt.expectedBody["message"].(string))).Once()

			}else {
				data = append(data, tt.expectedBody["public_song_lyrics"].([]models.PublicSongLyricResponse)...)
				tt.expectedBody["next"] = utils.GenerateNextLink(c, len(data), url.Values{
					"term": {term},
					"offset": {strconv.Itoa(offsetInt + 5)},
				}.Encode())

				mockPubSongLyricsRepo.On("SearchSongLyricsByTermShazam", term, "artists,songs", "5", offsetInt).Return(data, nil).Once()
			}

			err = pubSongLyricsController.SearchTermSongLyrics(c)

			if tt.wantErr {
				assert.Error(t, err)
				httpErr, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedCode, httpErr.Code)
				assert.Equal(t, tt.expectedBody, httpErr.Message)
			}else {
				assert.NoError(t, err)
				
				var ret struct {
					Next string `json:"next"`
					PublicSongLryics []models.PublicSongLyricResponse `json:"public_song_lyrics"`
				}
				err = json.Unmarshal(rec.Body.Bytes(), &ret)

				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, rec.Code)
				assert.Equal(t, tt.expectedBody["next"], ret.Next)
				assert.Equal(t, tt.expectedBody["public_song_lyrics"], ret.PublicSongLryics)

			}
			

		})
	}
}


func TestSearchAudioSongLyric(t *testing.T){
	mockPubSongLyricsRepo := mockService.MockPublicSongLyricsRepo{}	
	services.SetPublicSongLyricsRepo(&mockPubSongLyricsRepo)
	pubSongLyricsController := NewPublicSongLyricsController(&mockPubSongLyricsRepo)

	tests := []struct {
		name         string
		pathFile 		 string
		expectedBody echo.Map
		expectedCode int
		wantErr      bool
	}{
		{
			name: "Success",
			pathFile: "./test_file_audio/success.mp3",
			expectedCode: http.StatusOK,
			expectedBody: echo.Map{
				"public_song_lyrics": models.PublicSongLyricResponse{
					Title: "test",
					ArtistNames: "test",
					Lyric: "test",
				},
			},
			wantErr:      false,
		},
		{
			name: "Failed: not audio file",
			pathFile: "./test_file_audio/not_audio_file.txt",
			expectedCode: http.StatusBadRequest,
			expectedBody: echo.Map{
				"message": "invalid file type. please upload an audio file",
			},
			wantErr:      true,
		},
		{
			name: "Failed: size more than 500kb",
			pathFile: "./test_file_audio/more_than_500kb.mp3",
			expectedCode: http.StatusRequestEntityTooLarge,
			expectedBody: echo.Map{
				"message": "audio size must be less than 500kb",
			},
			wantErr:      true,
		},
		{
			name: "Failed: song lyric not found ",
			pathFile: "./test_file_audio/not_found.mp3",
			expectedCode: http.StatusNotFound,
			expectedBody: echo.Map{
				"message": "song lyric not found",
			},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		e := echo.New()
		
		// setup body for audio formfile audio
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("audio", tt.pathFile)
		assert.NoError(t, err)
		dataSended, err := os.Open(tt.pathFile)
		assert.NoError(t, err)
		_, err = io.Copy(part, dataSended)
		assert.NoError(t, err)
    assert.NoError(t, writer.Close())

		req := httptest.NewRequest(http.MethodGet, "/song_lyrics/public/search/audio", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType()) 
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		
		audioData, _ := c.FormFile("audio")
		rawBase64 := utils.Audio2RawBase64(audioData)

		if tt.name == "Failed: not audio file" {
			isAudio := utils.CheckAudioFile(audioData)
			assert.False(t, isAudio)
			assert.Empty(t, rawBase64)
		}else {
			assert.NotEmpty(t, rawBase64)
		}
		
		var data models.PublicSongLyricResponse

		if tt.wantErr {
			if tt.name != "Failed: not audio file" {
				mockPubSongLyricsRepo.On("SearchSongLyricByAudioRapidShazam", rawBase64).Return(data, errors.New(tt.expectedBody["message"].(string))).Once()
			}
		}else {
			data = tt.expectedBody["public_song_lyrics"].(models.PublicSongLyricResponse)
			mockPubSongLyricsRepo.On("SearchSongLyricByAudioRapidShazam", rawBase64).Return(data, nil).Once()
		}

		err = pubSongLyricsController.SearchAudioSongLyric(c)
		if tt.wantErr {
			assert.Error(t, err)
			httpErr, ok := err.(*echo.HTTPError)			
			assert.True(t, ok)
			assert.Equal(t, tt.expectedCode, httpErr.Code)
			assert.Equal(t, tt.expectedBody, httpErr.Message)
		}else {
			assert.NoError(t, err)
			
			var ret struct {
				PublicSongLryics models.PublicSongLyricResponse `json:"public_song_lyrics"`
			}
			err = json.Unmarshal(rec.Body.Bytes(), &ret)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)
			assert.Equal(t, tt.expectedBody["public_song_lyrics"], ret.PublicSongLryics)

		}


	}
	
}
