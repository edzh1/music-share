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
	apiToken:    "",
}

func (p *spotifyProvider) GetEntityID(URL, entity string) (string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}

	return u.Path, nil
}

func (p *spotifyProvider) GetName() string {
	return p.Name
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

	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		log.Println(resp.StatusCode)
		log.Fatal(string(b))
	}

	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		log.Fatal(err)
	}

	p.apiToken = fmt.Sprintf("%s %s", result.TokenType, result.AccessToken)
	log.Println(p.apiToken)

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
		if resp.StatusCode == 400 {
			p.Auth()
			request.Header.Set("authorization", p.apiToken)
			resp, err = client.Do(request)
		} else {
			b, _ := ioutil.ReadAll(resp.Body)
			log.Fatal(string(b))
		}
	}

	defer resp.Body.Close()

	var result getTrackResult

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return getTrackResult{}, err
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
		if resp.StatusCode == 400 {
			p.Auth()
			request.Header.Set("authorization", p.apiToken)
			resp, err = client.Do(request)
		} else {
			b, _ := ioutil.ReadAll(resp.Body)
			log.Fatal(string(b))
		}
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
		if resp.StatusCode == 400 {
			p.Auth()
			request.Header.Set("authorization", p.apiToken)
			resp, err = client.Do(request)
		}

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

func (p *spotifyProvider) Search(name, searchType string) (string, error) {
	query := url.QueryEscape(name)
	searchURL := fmt.Sprintf("%s?q=%s&type=%s", p.endpoints["SEARCH"], query, searchType)

	request, err := http.NewRequest("GET", searchURL, nil)
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

	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		log.Fatal(string(b))
	}

	defer resp.Body.Close()

	var result struct {
		Tracks struct {
			Items []struct {
				ID string
			}
		}
	}

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return result.Tracks.Items[0].ID, nil
}
