package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	// "github.com/edzh1/music-share/pkg/models"
)

type spotifyProvider struct {
	Provider
	ClientToken string
	apiToken    string
}

//Spotify provider
var Spotify = &spotifyProvider{
	Provider: Provider{
		Name: "spotify",
		endpoints: map[string]string{
			"GET_TRACKS":  "https://api.spotify.com/v1/tracks",
			"GET_ALBUMS":  "https://api.spotify.com/v1/albums",
			"GET_ARTISTS": "https://api.spotify.com/v1/artists",
			"SEARCH":      "https://api.spotify.com/v1/search",
			"AUTH":        "https://accounts.spotify.com/api/token",
		},
	},
	ClientToken: "",
	apiToken:    "Bearer BQBwP7xRtjPj5ZluO4g7w1n17GxQO2X4gkTDEnDvfESOsRFB6H5MuieEGyTQbUUu1VvUiutNimmmZ5XYxL8",
}

// type getSpotifyTrackResult struct {
// 	ID      string
// 	Name    string
// 	Artists []*struct {
// 		ID   string
// 		Name string
// 	}
// }

// type getSpotifyAlbumResult struct {
// 	ID   string
// 	Name string
// }

// type getSpotifyArtistResult struct {
// 	ID   string `json:"id"`
// 	Name string `json:"name"`
// }

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

	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		log.Fatal(string(b))
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

func (p *spotifyProvider) GetTrack(trackID string) (getTrackResult, error) {
	url := fmt.Sprintf("%s/%s", p.endpoints["GET_TRACKS"], trackID)
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("authorization", p.apiToken)

	if err != nil {
		log.Fatal(err)
		return getTrackResult{}, err
	}

	resp, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
		return getTrackResult{}, err
	}

	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		log.Fatal(string(b))
	}

	defer resp.Body.Close()

	var result getTrackResult

	err = json.NewDecoder(resp.Body).Decode(&result)

	// log.Println(result)

	if err != nil {
		log.Fatal(err)
		return getTrackResult(result), err
	}

	return result, nil
}

func (p *spotifyProvider) GetAlbum(albumID string) (getAlbumResult, error) {
	url := fmt.Sprintf("%s/%s", p.endpoints["GET_ALBUMS"], albumID)
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("authorization", p.apiToken)

	if err != nil {
		log.Fatal(err)
		return getAlbumResult{}, err
	}

	resp, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
		return getAlbumResult{}, err
	}

	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		log.Fatal(string(b))
	}

	defer resp.Body.Close()

	var result getAlbumResult

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		log.Fatal(err)
		return getAlbumResult{}, err
	}

	return result, nil
}

func (p *spotifyProvider) GetArtist(artistID string) (getArtistResult, error) {
	url := fmt.Sprintf("%s/%s", p.endpoints["GET_ARTISTS"], artistID)
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("authorization", p.apiToken)

	if err != nil {
		log.Fatal(err)
		return getArtistResult{}, err
	}

	resp, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
		return getArtistResult{}, err
	}

	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		log.Fatal(string(b))
	}

	defer resp.Body.Close()

	var result getArtistResult

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		log.Fatal(err)
		return getArtistResult{}, err
	}

	return getArtistResult(result), nil
}
