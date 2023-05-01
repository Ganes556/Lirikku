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

func RequestShazamMetadata(term, offset, types, limit string) (models.ShazamMetadata, error) {

	baseShazamMetadata := "https://" + os.Getenv("SHAZAM_API_HOST") + "/services/search/v4/id/ID/web/search"

	res, err := request(baseShazamMetadata, url.Values{
		"term": {term},
		"offset": {offset},
		"types": {types},
		"limit": {limit},
	}.Encode(), 10)

	if err != nil {
		return models.ShazamMetadata{}, err
	}

	defer res.Body.Close()

	var shazamMetadata models.ShazamMetadata

	json.NewDecoder(res.Body).Decode(&shazamMetadata)

	return shazamMetadata, nil

}

func RequestShazamLyric(key string) (models.ShazamLyric, error) {
	uriShazamLyric := "https://" + os.Getenv("SHAZAM_API_HOST") + "/discovery/v5/id/ID/web/-/track/" + key		
	
	res, _ := request(uriShazamLyric, url.Values{}.Encode(), 10)

	defer res.Body.Close()

	var shazamLyric models.ShazamLyric

	json.NewDecoder(res.Body).Decode(&shazamLyric)

	return shazamLyric, nil
}
