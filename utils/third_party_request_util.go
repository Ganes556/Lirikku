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

func request(method, url, param, body string, timeout int) (*http.Response, error){

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	var req *http.Request

	if method == "POST" {
		req, _ = http.NewRequest("POST", url, strings.NewReader(body))
	}else {
		req, _ = http.NewRequest("GET", url + "?" + param, nil)
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func RequestShazamSearchTerm(term, offset, types, limit string) (models.ReponseShazamSearchTerm, error) {

	baseShazamSearchTerm := "https://" + os.Getenv("SHAZAM_API_HOST") + "/services/search/v4/id/ID/web/search"

	res, err := request("GET",baseShazamSearchTerm, url.Values{
		"term": {term},
		"offset": {offset},
		"types": {types},
		"limit": {limit},
	}.Encode(),"", 10)

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
	
	res, _ := request("GET", urlShazamSearchKey, url.Values{}.Encode(),"", 10)

	defer res.Body.Close()

	var resData models.ResponseShazamSearchKey

	json.NewDecoder(res.Body).Decode(&resData)

	return resData, nil
}