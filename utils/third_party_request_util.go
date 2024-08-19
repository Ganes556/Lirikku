package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Lirikku/models"
)

func RequestTermByShazam(term, types string, offset, pageSize int) (models.ShazamSearchTermResponse, error) {

	urlShazamSearchKey := "https://www.shazam.com/services/search/v4/id/ID/web/search"

	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	query := url.Values{
		"term":   {term},
		"offset": {strconv.Itoa(offset)},
		"types":  {types},
		"limit":  {strconv.Itoa(pageSize)},
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

func RequestTermByOvh(term string) (models.OvhSearchTermResponse, error) {

	url := "https://api.lyrics.ovh/suggest/" + term

	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}
	
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return models.OvhSearchTermResponse{}, err
	}

	res, err := client.Do(req)

	if err != nil {
		return models.OvhSearchTermResponse{}, err
	}

	defer res.Body.Close()

	var resData models.OvhSearchTermResponse

	json.NewDecoder(res.Body).Decode(&resData)

	return resData, nil
}

func RequestLyricByOvh(artist string, title string) (models.OvhSearchLyricResponse, error) {
	url := fmt.Sprintf("https://api.lyrics.ovh/v1/%s/%s", artist, title)
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return models.OvhSearchLyricResponse{}, err
	}

	res, err := client.Do(req)

	if err != nil {
		return models.OvhSearchLyricResponse{}, err
	}

	defer res.Body.Close()
	var resData models.OvhSearchLyricResponse

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
		Timeout: time.Duration(50) * time.Second,
	}

	req, err := http.NewRequest("POST", urlShazamSearchAudio, strings.NewReader(rawBase64))
	req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("X-RapidAPI-Key", os.Getenv("RAPID_SHAZAM_API_KEY"))

	if err != nil {
		fmt.Println("err1->",err)
		return models.RapidShazamSearchAudioResponse{}, err
	}

	res, err := client.Do(req)

	if err != nil {
		fmt.Println("err2->",err)
		return models.RapidShazamSearchAudioResponse{}, err
	}

	defer res.Body.Close()

	var resData models.RapidShazamSearchAudioResponse
	// var resss map[string]any

	// json.NewDecoder(res.Body).Decode(&resss)

	// indent, _ := json.MarshalIndent(resss, "", " ")

	// fmt.Println(string(indent))

	json.NewDecoder(res.Body).Decode(&resData)

	fmt.Println(resData)

	return resData, nil
}
