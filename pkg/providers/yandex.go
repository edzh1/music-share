package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type yandexProvider struct {
	Provider
}

//Yandex provider
var Yandex = &yandexProvider{
	Provider: Provider{
		Name: "yandex",
		endpoints: map[string]string{
			"GET_TRACKS":  "https://music.yandex.ru/handlers/track.jsx",
			"GET_ALBUMS":  "https://music.yandex.ru/handlers/album.jsx",
			"GET_ARTISTS": "https://music.yandex.ru/handlers/artist.jsx",
			"SEARCH":      "https://music.yandex.ru/handlers/music-search.jsx",
		},
	},
}

type getYandexAlbumResult struct {
	ID   int
	Name string `json:"title"`
}

func (p *yandexProvider) GetName() string {
	return p.Name
}

func (p *yandexProvider) GetTrack(trackID string) (getTrackResult, error) {
	url := fmt.Sprintf("%s?track=%s&lang=en", p.endpoints["GET_TRACKS"], trackID)
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
		return getTrackResult{}, err
	}

	resp, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
		return getTrackResult{}, err
	}

	defer resp.Body.Close()

	var result struct {
		Track struct {
			ID    string
			Title string
		}
		Artists []*struct {
			ID   string
			Name string
		} `json:"artists"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		log.Fatal(err)
		return getTrackResult{}, err
	}

	return getTrackResult{
		ID:      result.Track.ID,
		Name:    result.Track.Title,
		Artists: result.Artists,
	}, nil
}

func (p *yandexProvider) GetAlbum(albumID string) (getAlbumResult, error) {
	url := fmt.Sprintf("%s?album=%s&lang=en", p.endpoints["GET_ALBUMS"], albumID)
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
		return getAlbumResult{}, err
	}

	resp, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
		return getAlbumResult{}, err
	}

	defer resp.Body.Close()

	var result getYandexAlbumResult

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		log.Fatal(err)
		return getAlbumResult{}, err
	}

	return getAlbumResult{
		ID:   strconv.Itoa(result.ID),
		Name: result.Name,
	}, nil
}

func (p *yandexProvider) GetArtist(artistID string) (getArtistResult, error) {
	url := fmt.Sprintf("%s?artist=%s&lang=en", p.endpoints["GET_ARTISTS"], artistID)
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
		return getArtistResult{}, err
	}

	resp, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
		return getArtistResult{}, err
	}

	defer resp.Body.Close()

	var result struct {
		Artist struct {
			ID   string
			Name string
		}
	}

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		log.Fatal(err)
		return getArtistResult{}, err
	}

	return getArtistResult{
		ID:   result.Artist.ID,
		Name: result.Artist.Name,
	}, err
}

func (p *yandexProvider) Search(name, searchType string) (string, error) {
	query := url.QueryEscape(name)
	searchURL := fmt.Sprintf("%s?text=%s&type=%s&lang=en", p.endpoints["SEARCH"], query, searchType)

	request, err := http.NewRequest("GET", searchURL, nil)

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
				ID int
			}
		}
	}

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return strconv.Itoa(result.Tracks.Items[0].ID), nil
}
