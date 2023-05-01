package utils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Lirikku/models"
)


func RequestShazamSearchTerm(term, offset, types, limit string) (models.ReponseShazamSearchTerm, error) {

	urlShazamSearchKey := "https://" + os.Getenv("SHAZAM_API_HOST") + "/services/search/v4/id/ID/web/search"

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
		return models.ReponseShazamSearchTerm{}, err
	}

	res, err := client.Do(req)


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
	
	client := &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}
	
	req, err := http.NewRequest("GET", urlShazamSearchKey, nil)

	if err != nil {
		return models.ResponseShazamSearchKey{}, err
	}

	res, err := client.Do(req)

	if err != nil {
		return models.ResponseShazamSearchKey{}, err
	}

	defer res.Body.Close()

	var resData models.ResponseShazamSearchKey

	json.NewDecoder(res.Body).Decode(&resData)

	return resData, nil
}