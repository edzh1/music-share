package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type spotifyProvider struct {
	provider
}

//Spotify provider
var Spotify = &spotifyProvider{
	provider: provider{
		Name:        "spotify",
		ClientToken: "",
		apiToken:    "",
		endpoints: map[string]string{
			"GET_TRACK": "https://api.spotify.com/v1/tracks",
			"GET_ALBUM": "https://api.spotify.com/v1/albums",
			"SEARCH":    "https://api.spotify.com/v1/search",
			"AUTH":      "https://accounts.spotify.com/api/token",
		},
	},
}

func (p *spotifyProvider) Auth() string {
	requestBody := url.Values{}
	requestBody.Set("grant_type", "client_credentials")

	url := Spotify.endpoints["AUTH"]
	request, err := http.NewRequest("POST", url, strings.NewReader(requestBody.Encode()))
	request.Header.Set("authorization", fmt.Sprintf("Basic %s", p.ClientToken))
	request.Header.Set("content-type", "application/x-www-form-urlencoded")

	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		log.Fatal(err)
	}

	p.apiToken = result.AccessToken

	return result.AccessToken
}

func (p *spotifyProvider) GetTrack(trackID string) string {
	url := fmt.Sprintf("%s/%s", Spotify.endpoints["GET_TRACK"], trackID)
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("authorization", p.apiToken)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}

func (p *spotifyProvider) GetAlbum(albumID string) string {
	url := fmt.Sprintf("%s/%s", Spotify.endpoints["GET_ALBUM"], albumID)
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("authorization", p.apiToken)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}
