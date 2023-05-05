package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Lirikku/middlewares"
	"github.com/Lirikku/models"
	"github.com/Lirikku/services"
	mockService "github.com/Lirikku/services/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {

	mockAuthRepo := mockService.MockAuthRepo{}
	services.SetAuthRepo(&mockAuthRepo)
	authController := NewAuthController(&mockAuthRepo)

	tests := []struct {
		name         string
		payload      models.UserRegister
		expectedBody echo.Map
		expectedCode int
		wantErr      bool
	}{
		{
			name: "Success",
			payload: models.UserRegister{
				Name: 	 "test",
				Email:    "test@gmail.com",
				Password: "test12345",
			},
			expectedBody: echo.Map{
				"message": "success created user",
			},
			expectedCode: http.StatusCreated,
			wantErr:      false,
		},
		{
			name: "Failed: email already registered",
			payload: models.UserRegister{
				Name: 	 "test",
				Email:    "test@gmail.com",
				Password: "test12345",
			},
			expectedBody: echo.Map{
				"message": "email already registered",
			},
			expectedCode: http.StatusBadRequest,
			wantErr:      true,
		},
		{
			name: "Failed: invalid email",
			payload: models.UserRegister{
				Name: "test2",
				Email: "test2gmail.com",
				Password: "test2123456",
			},
			expectedBody: echo.Map{
				"message": "email is not valid",
			},
			expectedCode: http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "Failed: invalid password",
			payload: models.UserRegister{
				Name: "test2",
				Email: "test2@gmail.com",
				Password: "test",
			},
			expectedBody: echo.Map{
				"message": "password must be at least 8 characters",
			},
			expectedCode: http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "Failed: required name",
			payload: models.UserRegister{
				Email: "test2@gmail.com",
				Password: "test2123456",
			},
			expectedBody: echo.Map{
				"message": "name is required",
			},
			expectedCode: http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "Failed: required email",
			payload: models.UserRegister{
				Name: "test2",
				Password: "test2123456",
			},
			expectedBody: echo.Map{
				"message": "email is required",
			},
			expectedCode: http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "Failed: required password",
			payload: models.UserRegister{
				Name: "test2",
				Email: "test2@gmail.com",
			},
			expectedBody: echo.Map{
				"message": "password is required",
			},
			expectedCode: http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "Internal Error: create user",
			payload: models.UserRegister{
				Name: "test2",
				Email: "testing@gmail.com",
				Password: "test2123456",
			},
			expectedBody: echo.Map{
				"message": "internal server error",
			},
			expectedCode: http.StatusInternalServerError,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = middlewares.NewValidator()
			
			payload, err := json.Marshal(tt.payload)
			assert.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			
			if tt.wantErr {
				
				if tt.name == "Internal Error: create user" {
					mockAuthRepo.On("CheckUserEmail", tt.payload.Email).Return(nil).Once()
					mockAuthRepo.On("CreateUser", tt.payload).Return(errors.New("internal server error")).Once()
				}else {

					mockAuthRepo.On("CheckUserEmail", tt.payload.Email).Return(errors.New(tt.expectedBody["message"].(string))).Once()
				}

			}else {
				mockAuthRepo.On("CheckUserEmail", tt.payload.Email).Return(nil).Once()
				mockAuthRepo.On("CreateUser", tt.payload).Return(nil).Once()
			}

			c := e.NewContext(req, rec)

			err = authController.Register(c)
			
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

func TestLogin(t *testing.T) {
	mockAuthRepo := mockService.MockAuthRepo{}
	services.SetAuthRepo(&mockAuthRepo)
	authController := NewAuthController(&mockAuthRepo)

	tests := []struct {
		name         string
		payload      models.UserLogin
		expectedBody echo.Map
		expectedCode int
		wantErr      bool
	}{
		{
			name: "Success",
			payload: models.UserLogin{
				Email:    "test@gmail.com",
				Password: "test12345",
			},
			expectedBody: echo.Map{
				"message": "success login",
				"token": true,
			},
			expectedCode: http.StatusOK,
			wantErr:      false,
		},
		{
			name: "Failed: email not registered",
			payload: models.UserLogin{
				Email:    "test@gmail.com",
				Password: "test12345",
			},
			expectedBody: echo.Map{
				"message": "email not registered",
			},
			expectedCode: http.StatusUnauthorized,
			wantErr:      true,
		},
		{
			name: "Failed: login failed",
			payload: models.UserLogin{
				Email: "test@gmail.com",
				Password: "test2123456",
			},
			expectedBody: echo.Map{
				"message": "incorrect email or password",
			},
			expectedCode: http.StatusUnauthorized,
			wantErr: true,
		},
		{
			name: "Failed: invalid email",
			payload: models.UserLogin{
				Email: "test2gmail.com",
				Password: "test2123456",
			},
			expectedBody: echo.Map{
				"message": "email is not valid",
			},
			expectedCode: http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "Failed: invalid password",
			payload: models.UserLogin{
				Email: "test2@gmail.com",
				Password: "test",
			},
			expectedBody: echo.Map{
				"message": "password must be at least 8 characters",
			},
			expectedCode: http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "Failed: required email",
			payload: models.UserLogin{
				Password: "test2123456",
			},
			expectedBody: echo.Map{
				"message": "email is required",
			},
			expectedCode: http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "Failed: required password",
			payload: models.UserLogin{
				Email: "test2@gmail.com",
			},
			expectedBody: echo.Map{
				"message": "password is required",
			},
			expectedCode: http.StatusBadRequest,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Validator = middlewares.NewValidator()
			
			payload, err := json.Marshal(tt.payload)
			assert.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			
			var data models.User

			if tt.wantErr {
				if tt.expectedBody["message"] == "email not registered" {
					mockAuthRepo.On("GetUserByEmail", tt.payload.Email).Return(models.User{},errors.New(tt.expectedBody["message"].(string))).Once()
				}
				if tt.expectedBody["message"] == "incorrect email or password" {
					data = models.User{
						Base: models.Base{
							ID: 1,
						},
						Name: "test",
						Email: tt.payload.Email,
						Password: "$2y$10$kNimIla567HkdVcawxIPfu6gC0KCC2mpibPUnxfxobBhcMIODAI.K",
					}
					mockAuthRepo.On("GetUserByEmail", tt.payload.Email).Return(data, nil).Once()
				}
			}else {
				data = models.User{
					Base: models.Base{
						ID: 1,
					},
					Name: "test",
					Email: tt.payload.Email,
					Password: "$2y$10$kNimIla567HkdVcawxIPfu6gC0KCC2mpibPUnxfxobBhcMIODAI.K",
				}
				mockAuthRepo.On("GetUserByEmail", tt.payload.Email).Return(data, nil).Once()
			}

			c := e.NewContext(req, rec)

			err = authController.Login(c)

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
				assert.Equal(t, tt.expectedBody["message"], ret["message"])
				if tt.expectedBody["token"] == true {
					assert.NotEmpty(t, ret["token"])
				}

			}

		})

	}
}