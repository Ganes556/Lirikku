package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/Lirikku/middlewares"
	"github.com/Lirikku/models"
	"github.com/Lirikku/services"
	mockService "github.com/Lirikku/services/mocks"
	"github.com/Lirikku/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)


func TestGetSongLyrics(t *testing.T){
	mockMySongLyricsRepo := mockService.MockMySongLyricsRepo{}	
	services.SetMySongLyricsRepo(&mockMySongLyricsRepo)
	mySongLyricsController := NewMySongLyricsController(&mockMySongLyricsRepo)


	tests := []struct {
		name         string
		param 			 string
		expectedBody echo.Map
		expectedCode int
		wantErr      bool
	}{
		{
			name: "Success",
			param: "?offset=1",
			expectedCode: http.StatusOK,
			expectedBody: echo.Map{
				"next": "",
				"my_song_lyrics": []models.SongLyricResponse{
					{
						ID: 1,
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
			param: "?offset=abc",
			expectedCode: http.StatusBadRequest,
			expectedBody: echo.Map{
				"message": "offset must be a number",
			},
			wantErr:      true,
		},
		{
			name: "Failed: internal server error (GetSongLyrics)",
			param: "?offset=1",
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

			req := httptest.NewRequest(http.MethodGet, "/song_lyrics/my" + tt.param, nil)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			
			c.Set("user", models.UserJWTDecode{
				ID: 1,
				Name: "test",
			})

			offset := c.QueryParam("offset")
			
			offsetInt, err := utils.CheckOffset(offset)

			if tt.name != "Failed: offset must be a number" {
				assert.NoError(t, err)
			}

			var data []models.SongLyricResponse
			if tt.wantErr {
				mockMySongLyricsRepo.On("GetSongLyrics", uint(1), offsetInt).Return(data,errors.New("internal server error")).Once()
			}else {
				data = append(data, models.SongLyricResponse{
					ID: 1,
					Title: "test",
					ArtistNames: "test",
					Lyric: "test",
				})

				tt.expectedBody["next"] = utils.GenerateNextLink(c, len(data), url.Values{
					"offset": {strconv.Itoa(offsetInt + 5)},
				}.Encode())
				
				mockMySongLyricsRepo.On("GetSongLyrics", uint(1), offsetInt).Return(data, nil).Once()

			}

			err = mySongLyricsController.GetSongLyrics(c)
			
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
					MySongLyrics []models.SongLyricResponse `json:"my_song_lyrics"`
				}
				err = json.Unmarshal(rec.Body.Bytes(), &ret)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, rec.Code)
				assert.Equal(t, tt.expectedBody["next"], ret.Next)
				assert.Equal(t, tt.expectedBody["my_song_lyrics"], ret.MySongLyrics)
			}

		})
	}
}
func TestGetSongLyric(t *testing.T){

	mockMySongLyricsRepo := mockService.MockMySongLyricsRepo{}	
	services.SetMySongLyricsRepo(&mockMySongLyricsRepo)
	mySongLyricsController := NewMySongLyricsController(&mockMySongLyricsRepo)
	
	tests := []struct {
		name         string
		idSongLyric  string
		expectedBody echo.Map
		expectedCode int
		wantErr      bool
	}{
		{
			name: "Success",
			idSongLyric: "1",
			expectedCode: http.StatusOK,
			expectedBody: echo.Map{
				"my_song_lyrics": models.SongLyricResponse{
					ID: 1,
					Title: "test",
					ArtistNames: "test",
					Lyric: "test",
				},
			},
			wantErr:      false,
		},
		{
			name: "Failed: id must be a number",
			idSongLyric: "abc",
			expectedCode: http.StatusBadRequest,
			expectedBody: echo.Map{
				"message": "id must be a number",
			},
			wantErr:      true,
		},
		{
			name: "Failed: song lyric not found (GetSongLyric)",
			idSongLyric: "99",
			expectedCode: http.StatusBadRequest,
			expectedBody: echo.Map{
				"message": "song lyric not found",
			},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/song_lyrics/my/" + tt.idSongLyric, nil)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			user := models.UserJWTDecode{
				ID: 1,
				Name: "test",
			}

			c.Set("user", user)

			c.SetParamNames("id")
			c.SetParamValues(tt.idSongLyric)
			
			
			num, err := utils.CheckId(tt.idSongLyric)
			
			if tt.name != "Failed: id must be a number" {
				assert.NoError(t,err)
			}

			var data models.SongLyricResponse
			
			if tt.wantErr {
				mockMySongLyricsRepo.On("GetSongLyric", num , user.ID).Return(data, errors.New("song lyric not found")).Once()
			}else {
				data = tt.expectedBody["my_song_lyrics"].(models.SongLyricResponse)
				mockMySongLyricsRepo.On("GetSongLyric", 1, user.ID).Return(data, nil).Once()
			}

			err = mySongLyricsController.GetSongLyric(c)
			
			if tt.wantErr {
				assert.Error(t, err)
				httpErr, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedCode, httpErr.Code)
				assert.Equal(t, tt.expectedBody, httpErr.Message)

			}else {
				assert.NoError(t, err)

				var ret struct {
					MySongLyrics models.SongLyricResponse `json:"my_song_lyrics"`
				}
				err = json.Unmarshal(rec.Body.Bytes(), &ret)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, rec.Code)
			}

		})
	}
}

func TestSaveSongLyric(t *testing.T){
	mockMySongLyricsRepo := mockService.MockMySongLyricsRepo{}	
	services.SetMySongLyricsRepo(&mockMySongLyricsRepo)
	mySongLyricsController := NewMySongLyricsController(&mockMySongLyricsRepo)

	tests := []struct {
		name         string
		payload 		 models.SongLyricWrite
		expectedBody echo.Map
		expectedCode int
		wantErr      bool
	}{
		{
			name: "Success",
			payload: models.SongLyricWrite{
				Title: "test",
				ArtistNames: "test",
				Lyric: "test",
			},
			expectedCode: http.StatusCreated,
			expectedBody: echo.Map{
				"message": "song lyric saved successfully",
			},
			wantErr:      false,
		},
		{
			name: "Failed: internal server error (SaveSongLyric)",
			payload: models.SongLyricWrite{
				Title: "test",
				ArtistNames: "test",
				Lyric: "test",
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: echo.Map{
				"message": "internal server error",
			},
			wantErr:      true,
		},
		{
			name: "Failed: song lyric already saved",
			payload: models.SongLyricWrite{
				Title: "test",
				ArtistNames: "test",
				Lyric: "test",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: echo.Map{
				"message": "song lyric already saved",
			},
			wantErr:      true,
		},
		{
			name: "Failed: title required",
			payload: models.SongLyricWrite{
				ArtistNames: "test",
				Lyric: "test",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: echo.Map{
				"message": "title is required",
			},
			wantErr:      true,
		},
		{
			name: "Failed: artist_names required",
			payload: models.SongLyricWrite{
				Title: "test",
				Lyric: "test",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: echo.Map{
				"message": "artist_names is required",
			},
			wantErr:      true,
		},
		{
			name: "Failed: lyric required",
			payload: models.SongLyricWrite{
				Title: "test",
				ArtistNames: "test",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: echo.Map{
				"message": "lyric is required",
			},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = middlewares.NewValidator()
			
			payload, err := json.Marshal(tt.payload)
			assert.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, "/song_lyrics/my", bytes.NewReader(payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			
			c := e.NewContext(req, rec)

			user := models.UserJWTDecode{
				ID: 1,
				Name: "test",
			}

			c.Set("user", user)
			
			if tt.wantErr {
				if tt.name != "Failed: internal server error (SaveSongLyric)" {
					mockMySongLyricsRepo.On("CheckSongLyric", user.ID, tt.payload).Return(errors.New(tt.expectedBody["message"].(string))).Once()
				}else {
					mockMySongLyricsRepo.On("CheckSongLyric", user.ID, tt.payload).Return(nil).Once()
					mockMySongLyricsRepo.On("SaveSongLyric", user.ID, tt.payload).Return(errors.New(tt.expectedBody["message"].(string))).Once()
				}
			}else {
				mockMySongLyricsRepo.On("CheckSongLyric", user.ID, tt.payload).Return(nil).Once()
				mockMySongLyricsRepo.On("SaveSongLyric", user.ID, tt.payload).Return(nil).Once()
			}

			err = mySongLyricsController.SaveSongLyric(c)
			
			if tt.wantErr {
				assert.Error(t, err)
				httpErr, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedCode, httpErr.Code)
				assert.Equal(t, tt.expectedBody, httpErr.Message)

			}else {
				assert.NoError(t, err)
				var ret echo.Map
				err = json.Unmarshal(rec.Body.Bytes(), &ret)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, rec.Code)
				assert.Equal(t, tt.expectedBody, ret)
			}

		})
	}
}

func TestSearchSongLyrics(t *testing.T){

	mockMySongLyricsRepo := mockService.MockMySongLyricsRepo{}	
	services.SetMySongLyricsRepo(&mockMySongLyricsRepo)
	mySongLyricsController := NewMySongLyricsController(&mockMySongLyricsRepo)

	tests := []struct {
		name         string
		param 		 string
		expectedBody echo.Map
		expectedCode int
		wantErr      bool
	}{
		{
			name: "Success",
			param: "?title=test&artist_names=test&lyric=test",
			expectedCode: http.StatusOK,
			expectedBody: echo.Map{
				"next": "",
				"my_song_lyrics": []models.SongLyricResponse{
					{
						ID: 1,
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
			param: "?title=test&artist_names=test&lyric=test&offset=abc",
			expectedCode: http.StatusBadRequest,
			expectedBody: echo.Map{
				"message": "offset must be a number",
			},
			wantErr:      true,
		},
		{
			name: "Failed: internal server error (SearchSongLyrics)",
			param: "?title=test&artist_names=test&lyric=test",
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
			
			req := httptest.NewRequest(http.MethodPost, "/song_lyrics/my/search" + tt.param, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
		
			c := e.NewContext(req, rec)

			user := models.UserJWTDecode{
				ID: 1,
				Name: "test",
			}

			c.Set("user", user)
			title := c.QueryParam("title")
			lyric := c.QueryParam("lyric")
			artistNames := c.QueryParam("artist_names")
			
			offset := c.QueryParam("offset")

			offsetInt, err := utils.CheckOffset(offset)

			if tt.name != "Failed: offset must be a number" {
				assert.NoError(t, err)
			}
			
			var data []models.SongLyricResponse

			if tt.wantErr {
				mockMySongLyricsRepo.On("SearchSongLyrics", user.ID, title, lyric, artistNames, offsetInt).Return(data, errors.New(tt.expectedBody["message"].(string))).Once()
			}else {
				data = append(data, tt.expectedBody["my_song_lyrics"].([]models.SongLyricResponse)...)

				tt.expectedBody["next"] = utils.GenerateNextLink(c, len(data), url.Values{
					"title": {title},
					"lyric": {lyric},
					"artist_names": {artistNames},
					"offset": {strconv.Itoa(offsetInt + 5)},
				}.Encode())
				mockMySongLyricsRepo.On("SearchSongLyrics", user.ID, title, lyric, artistNames, offsetInt).Return(data, nil).Once()
			}

			err = mySongLyricsController.SearchSongLyrics(c)
			
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
					MySongLyrics []models.SongLyricResponse `json:"my_song_lyrics"`
				}
				err = json.Unmarshal(rec.Body.Bytes(), &ret)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, rec.Code)
				assert.Equal(t, tt.expectedBody["next"], ret.Next)
				assert.Equal(t, tt.expectedBody["my_song_lyrics"], ret.MySongLyrics)
			}

		})
	}
}

func TestDeleteSongLyric(t *testing.T) {
	mockMySongLyricsRepo := mockService.MockMySongLyricsRepo{}	
	services.SetMySongLyricsRepo(&mockMySongLyricsRepo)
	mySongLyricsController := NewMySongLyricsController(&mockMySongLyricsRepo)

	tests := []struct {
		name         string
		idSongLyric	 string
		expectedBody echo.Map
		expectedCode int
		wantErr      bool
	}{
		{
			name: "Success",
			idSongLyric: "1",
			expectedCode: http.StatusOK,
			expectedBody: echo.Map{
				"message": "song lyric deleted successfully",
			},
			wantErr:      false,
		},
		{
			name: "Failed: id must be a number",
			expectedCode: http.StatusBadRequest,
			idSongLyric: "abc",
			expectedBody: echo.Map{
				"message": "id must be a number",
			},
			wantErr:      true,
		},
		{
			name: "Failed: song lyric not found",
			expectedCode: http.StatusBadRequest,
			idSongLyric: "123",
			expectedBody: echo.Map{
				"message": "song lyric not found",
			},
			wantErr:      true,
		},
		{
			name: "Failed: internal server error (DeleteSongLyrics)",
			idSongLyric: "1",
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
			
			req := httptest.NewRequest(http.MethodDelete, "/song_lyrics/my/" + tt.idSongLyric, nil)

			rec := httptest.NewRecorder()
		
			c := e.NewContext(req, rec)

			c.SetParamNames("id")
			c.SetParamValues(tt.idSongLyric)

			user := models.UserJWTDecode{
				ID: 1,
				Name: "test",
			}

			c.Set("user", user)
			
			idSongLyricInt, err := utils.CheckId(tt.idSongLyric)

			if tt.name != "Failed: id must be a number" {
				assert.NoError(t, err)
			}

			if tt.wantErr {
				if tt.name == "Failed: internal server error (DeleteSongLyrics)" {
					mockMySongLyricsRepo.On("GetSongLyric", idSongLyricInt, user.ID).Return(models.SongLyricResponse{}, nil).Once()
					mockMySongLyricsRepo.On("DeleteSongLyric", idSongLyricInt, user.ID).Return(errors.New(tt.expectedBody["message"].(string))).Once()
				}
				
				if tt.name == "Failed: song lyric not found" {
					mockMySongLyricsRepo.On("GetSongLyric", idSongLyricInt, user.ID).Return(models.SongLyricResponse{}, errors.New(tt.expectedBody["message"].(string))).Once()
				}

			}else {
				mockMySongLyricsRepo.On("GetSongLyric", idSongLyricInt, user.ID).Return(models.SongLyricResponse{}, nil).Once()
				mockMySongLyricsRepo.On("DeleteSongLyric", idSongLyricInt, user.ID).Return(nil).Once()
			}

			err = mySongLyricsController.DeleteSongLyric(c)
			
			if tt.wantErr {
				assert.Error(t, err)
				httpErr, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedCode, httpErr.Code)
				assert.Equal(t, tt.expectedBody, httpErr.Message)

			}else {
				assert.NoError(t, err)
				var ret echo.Map
				err = json.Unmarshal(rec.Body.Bytes(), &ret)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, rec.Code)
				assert.Equal(t, tt.expectedBody, ret)
			}

		})
	}

}

func TestUpdateSongLyric(t *testing.T) {
	mockMySongLyricsRepo := mockService.MockMySongLyricsRepo{}	
	services.SetMySongLyricsRepo(&mockMySongLyricsRepo)
	mySongLyricsController := NewMySongLyricsController(&mockMySongLyricsRepo)


	tests := []struct {
		name         string
		idSongLyric	 string
		payload 		 models.SongLyricWrite
		expectedBody echo.Map
		expectedCode int
		wantErr      bool
	}{
		{
			name: "Success",
			idSongLyric: "1",
			payload: models.SongLyricWrite{
				Title: "test",
				ArtistNames: "test",
			},
			expectedCode: http.StatusOK,
			expectedBody: echo.Map{
				"message": "song lyric updated successfully",
			},
			wantErr:      false,
		},
		{
			name: "Failed: id must be a number",
			idSongLyric: "abc",
			payload: models.SongLyricWrite{
				Title: "test",
				ArtistNames: "test",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: echo.Map{
				"message": "id must be a number",
			},
			wantErr:      true,
		},
		{
			name: "Failed: song lyric not found",
			idSongLyric: "123",
			payload: models.SongLyricWrite{
				Title: "test",
				ArtistNames: "test",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: echo.Map{
				"message": "song lyric not found",
			},
			wantErr:      true,
		},
		{
			name: "Failed: internal server error (UpdateSongLyric)",
			idSongLyric: "1",
			payload: models.SongLyricWrite{
				Title: "test",
				ArtistNames: "test",
			},
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
			
			payload, err := json.Marshal(tt.payload)
			assert.NoError(t, err)
			req := httptest.NewRequest(http.MethodPut, "/song_lyrics/my/" + tt.idSongLyric, bytes.NewBuffer(payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
		
			c := e.NewContext(req, rec)

			c.SetParamNames("id")
			c.SetParamValues(tt.idSongLyric)

			user := models.UserJWTDecode{
				ID: 1,
				Name: "test",
			}

			c.Set("user", user)
			
			idSongLyricInt, err := utils.CheckId(tt.idSongLyric)

			if tt.name != "Failed: id must be a number" {
				assert.NoError(t, err)
			}

			if tt.wantErr {
				if tt.name == "Failed: internal server error (UpdateSongLyric)" {
					mockMySongLyricsRepo.On("GetSongLyric", idSongLyricInt, user.ID).Return(models.SongLyricResponse{}, nil).Once()
					mockMySongLyricsRepo.On("UpdateSongLyric", idSongLyricInt, user.ID, tt.payload).Return(errors.New(tt.expectedBody["message"].(string))).Once()
				}
				
				if tt.name == "Failed: song lyric not found" {
					mockMySongLyricsRepo.On("GetSongLyric", idSongLyricInt, user.ID).Return(models.SongLyricResponse{}, errors.New(tt.expectedBody["message"].(string))).Once()
				}

			}else {
				mockMySongLyricsRepo.On("GetSongLyric", idSongLyricInt, user.ID).Return(models.SongLyricResponse{}, nil).Once()
				mockMySongLyricsRepo.On("UpdateSongLyric", idSongLyricInt, user.ID, tt.payload).Return(nil).Once()
			}

			err = mySongLyricsController.UpdateSongLyric(c)
			
			if tt.wantErr {
				assert.Error(t, err)
				httpErr, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedCode, httpErr.Code)
				assert.Equal(t, tt.expectedBody, httpErr.Message)

			}else {
				assert.NoError(t, err)
				var ret echo.Map
				err = json.Unmarshal(rec.Body.Bytes(), &ret)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCode, rec.Code)
				assert.Equal(t, tt.expectedBody, ret)
			}

		})
	}


}