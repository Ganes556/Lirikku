package utils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Lirikku/models"
)

func request(uri string, param string, timeout int) (*http.Response, error){

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	req, _ := http.NewRequest("GET", uri + "?" + param, nil)

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func RequestShazamSearchTerm(term, offset, types, limit string) (models.ReponseShazamSearchTerm, error) {

	baseShazamSearchTerm := "https://" + os.Getenv("SHAZAM_API_HOST") + "/services/search/v4/id/ID/web/search"

	res, err := request(baseShazamSearchTerm, url.Values{
		"term": {term},
		"offset": {offset},
		"types": {types},
		"limit": {limit},
	}.Encode(), 10)

	if err != nil {
		return models.ReponseShazamSearchTerm{}, err
	}

	defer res.Body.Close()

	var resData models.ReponseShazamSearchTerm

	json.NewDecoder(res.Body).Decode(&resData)

	return resData, nil

}

func RequestShazamSearchKey(key string) (models.ResponseShazamSearchKey, error) {
	urlShazamSearchKey := "https://" + os.Getenv("SHAZAM_API_HOST") + "/discovery/v5/id/ID/web/-/track/" + key		
	
	res, _ := request(urlShazamSearchKey, url.Values{}.Encode(), 10)

	defer res.Body.Close()

	var resData models.ResponseShazamSearchKey

	json.NewDecoder(res.Body).Decode(&resData)

	return resData, nil
}