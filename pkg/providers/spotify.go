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
	ClientToken string
	apiToken    string
}

//Spotify provider
var Spotify = &spotifyProvider{
	ClientToken: "",
	provider: provider{
		Name: "spotify",
		endpoints: map[string]string{
			"GET_TRACKS":  "https://api.spotify.com/v1/tracks",
			"GET_ALBUMS":  "https://api.spotify.com/v1/albums",
			"GET_ARTISTS": "https://api.spotify.com/v1/artists",
			"SEARCH":      "https://api.spotify.com/v1/search",
			"AUTH":        "https://accounts.spotify.com/api/token",
		},
	},
}

func (p *spotifyProvider) Auth() string {
	requestBody := url.Values{}
	requestBody.Set("grant_type", "client_credentials")

	url := p.endpoints["AUTH"]
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

func (p *spotifyProvider) GetTrack(trackID string) (string, error) {
	url := fmt.Sprintf("%s/%s", p.endpoints["GET_TRACKS"], trackID)
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("authorization", p.apiToken)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	resp, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return string(body), nil
}

func (p *spotifyProvider) GetAlbum(albumID string) (string, error) {
	url := fmt.Sprintf("%s/%s", p.endpoints["GET_ALBUMS"], albumID)
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("authorization", p.apiToken)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	resp, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return string(body), nil
}

func (p *spotifyProvider) GetArtist(artistID string) (string, error) {
	url := fmt.Sprintf("%s/%s", p.endpoints["GET_ARTISTS"], artistID)
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("authorization", p.apiToken)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	resp, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return string(body), nil
}
