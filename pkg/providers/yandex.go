package providers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

type getYandexTrackResult struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Artist []*struct {
		ID   string
		Name string
	} `json:"artists"`
}

type getYandexAlbumResult struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type getYandexArtistResult struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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
		Artist []*struct {
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
		Artists: result.Artist,
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

	return getAlbumResult(result), nil
}

func (p *yandexProvider) GetArtist(artistID string) (getArtistResult, error) {
	url := fmt.Sprintf("%s?artist=%s&lang=en", p.endpoints["GET_ALBUMS"], artistID)
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

	var result getYandexArtistResult

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		log.Fatal(err)
		return getArtistResult{}, err
	}

	return getArtistResult(result), err
}
