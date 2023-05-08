package utils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Lirikku/models"
)


func RequestShazamSearchTerm(term, offset, types, limit string) (models.ShazamSearchTermResponse, error) {

	urlShazamSearchKey := "https://www.shazam.com/services/search/v4/id/ID/web/search"

	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	query := url.Values{
		"term": {term},
		"offset": {offset},
		"types": {types},
		"limit": {limit},
	}.Encode()

	req, err := http.NewRequest("GET", urlShazamSearchKey+"?"+query, nil)

	if err != nil {
		return models.ShazamSearchTermResponse{}, err
	}

	res, err := client.Do(req)


	if err != nil {
		return models.ShazamSearchTermResponse{}, err
	}

	defer res.Body.Close()

	var resData models.ShazamSearchTermResponse

	json.NewDecoder(res.Body).Decode(&resData)

	return resData, nil

}

func RequestShazamSearchKey(key string) (models.ShazamSearchKeyResponse, error) {
	urlShazamSearchKey := "https://www.shazam.com/discovery/v5/id/ID/web/-/track/" + key		
	
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}
	
	req, err := http.NewRequest("GET", urlShazamSearchKey, nil)

	if err != nil {
		return models.ShazamSearchKeyResponse{}, err
	}

	res, err := client.Do(req)

	if err != nil {
		return models.ShazamSearchKeyResponse{}, err
	}

	defer res.Body.Close()

	var resData models.ShazamSearchKeyResponse

	json.NewDecoder(res.Body).Decode(&resData)

	return resData, nil
}

func RequestShazamSearchAudio(rawBase64 string) (models.RapidShazamSearchAudioResponse, error) {
	
	urlShazamSearchAudio := "https://shazam.p.rapidapi.com/songs/v2/detect"
	
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}
	
	req, err := http.NewRequest("POST", urlShazamSearchAudio, strings.NewReader(rawBase64))
	req.Header.Add("content-type", "text/plain")
	req.Header.Add("X-RapidAPI-Key", os.Getenv("RAPID_SHAZAM_API_KEY"))

	if err != nil {
		return models.RapidShazamSearchAudioResponse{}, err
	}
	
	res, err := client.Do(req)

	if err != nil {
		return models.RapidShazamSearchAudioResponse{}, err
	}
	
	defer res.Body.Close()

	var resData models.RapidShazamSearchAudioResponse

	json.NewDecoder(res.Body).Decode(&resData)
	
	return resData, nil
}
