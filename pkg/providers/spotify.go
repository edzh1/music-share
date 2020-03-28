package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type spotifyProvider struct {
	Provider
	ClientToken string
	apiToken    string
}

var ErrInvalidTokenMessage = "Only valid bearer authentication supported"
var ErrExpiredTokenMessage = "The access token expired"

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
		return "", ErrBadRequest
	}

	return path.Base(u.Path), nil
}

func (p *spotifyProvider) GetName() string {
	return p.Name
}

func (p *spotifyProvider) Auth() error {
	requestBody := url.Values{}
	requestBody.Set("grant_type", "client_credentials")

	url := p.endpoints["AUTH"]
	request, err := http.NewRequest("POST", url, strings.NewReader(requestBody.Encode()))
	request.Header.Set("authorization", fmt.Sprintf("Basic %s", p.ClientToken))
	request.Header.Set("content-type", "application/x-www-form-urlencoded")

	if err != nil {
		return ErrProviderFailure
	}

	resp, err := client.Do(request)

	if err != nil {
		return ErrProviderFailure
	}

	if resp.StatusCode != 200 {
		return ErrProviderFailure
	}

	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return ErrProviderFailure
	}

	p.apiToken = fmt.Sprintf("%s %s", result.TokenType, result.AccessToken)

	return nil
}

func (p *spotifyProvider) GetTrack(trackID string) (getTrackResult, error) {
	url := fmt.Sprintf("%s/%s", p.endpoints["GET_TRACKS"], trackID)
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("authorization", p.apiToken)

	if err != nil {
		return getTrackResult{}, ErrProviderFailure
	}

	resp, err := client.Do(request)

	if err != nil {
		return getTrackResult{}, ErrProviderFailure
	}

	err = p.handleError(resp)

	if err == ErrAuth {
		err = p.Auth()
		request.Header.Set("authorization", p.apiToken)
		resp, err = client.Do(request)

		err = p.handleError(resp)

		if err != nil {
			return getTrackResult{}, err
		}
	} else {
		return getTrackResult{}, err
	}

	defer resp.Body.Close()

	var result getTrackResult

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return getTrackResult{}, ErrProviderFailure
	}

	return result, nil
}

func (p *spotifyProvider) GetAlbum(albumID string) (getAlbumResult, error) {
	url := fmt.Sprintf("%s/%s", p.endpoints["GET_ALBUMS"], albumID)
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("authorization", p.apiToken)

	if err != nil {
		return getAlbumResult{}, ErrProviderFailure
	}

	resp, err := client.Do(request)

	if err != nil {
		return getAlbumResult{}, ErrProviderFailure
	}

	err = p.handleError(resp)

	if err == ErrAuth {
		err = p.Auth()
		request.Header.Set("authorization", p.apiToken)
		resp, err = client.Do(request)

		err = p.handleError(resp)

		if err != nil {
			return getAlbumResult{}, err
		}
	} else {
		return getAlbumResult{}, err
	}

	defer resp.Body.Close()

	var result getAlbumResult

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return getAlbumResult{}, ErrProviderFailure
	}

	return result, nil
}

func (p *spotifyProvider) GetArtist(artistID string) (getArtistResult, error) {
	url := fmt.Sprintf("%s/%s", p.endpoints["GET_ARTISTS"], artistID)
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("authorization", p.apiToken)

	if err != nil {
		return getArtistResult{}, ErrProviderFailure
	}

	resp, err := client.Do(request)

	if err != nil {
		return getArtistResult{}, ErrProviderFailure
	}

	err = p.handleError(resp)

	if err == ErrAuth {
		err = p.Auth()
		request.Header.Set("authorization", p.apiToken)
		resp, err = client.Do(request)

		err = p.handleError(resp)

		if err != nil {
			return getArtistResult{}, err
		}
	} else {
		return getArtistResult{}, err
	}

	defer resp.Body.Close()

	var result getArtistResult

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return getArtistResult{}, ErrProviderFailure
	}

	return getArtistResult(result), nil
}

func (p *spotifyProvider) Search(name, searchType string) (string, error) {
	query := url.QueryEscape(name)
	searchURL := fmt.Sprintf("%s?q=%s&type=%s", p.endpoints["SEARCH"], query, searchType)

	request, err := http.NewRequest("GET", searchURL, nil)
	request.Header.Set("authorization", p.apiToken)

	if err != nil {
		return "", ErrProviderFailure
	}

	resp, err := client.Do(request)

	if err != nil {
		return "", ErrProviderFailure
	}

	err = p.handleError(resp)

	if err == ErrAuth {
		err = p.Auth()
		request.Header.Set("authorization", p.apiToken)
		resp, err = client.Do(request)

		err = p.handleError(resp)

		if err != nil {
			return "", err
		}
	} else {
		return "", err
	}

	defer resp.Body.Close()

	var result struct {
		Tracks struct {
			Items []struct {
				ID string
			}
		}
		Albums struct {
			Items []struct {
				ID string
			}
		}
		Artists struct {
			Items []struct {
				ID string
			}
		}
	}

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return "", ErrProviderFailure
	}

	switch searchType {
	case "track":
		return (result.Tracks.Items[0].ID), nil
	case "album":
		return (result.Albums.Items[0].ID), nil
	case "artist":
		return (result.Artists.Items[0].ID), nil
	}

	return "", ErrWrongSearchType
}

func (p *spotifyProvider) handleError(resp *http.Response) error {
	var responseError struct {
		Error struct {
			Status  int
			Message string
		}
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 400 || resp.StatusCode == 401 {
			err := json.NewDecoder(resp.Body).Decode(&responseError)

			if err != nil {
				return ErrProviderFailure
			}

			if responseError.Error.Message == ErrInvalidTokenMessage ||
				responseError.Error.Message == ErrExpiredTokenMessage {
				return ErrAuth
			}

			return ErrBadRequest
		} else if resp.StatusCode == 404 {
			return ErrNotFound
		} else {
			return ErrProviderFailure
		}
	}

	return nil
}
