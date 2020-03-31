package providers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
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

func (p *yandexProvider) GetEntityID(URL, entity string) (string, error) {
	u, err := url.Parse(URL)

	if err != nil {
		return "", ErrBadRequest
	}

	return path.Base(u.Path), nil
}

func (p *yandexProvider) GetName() string {
	return p.Name
}

func (p *yandexProvider) GetTrack(trackID string) (getTrackResult, error) {
	url := fmt.Sprintf("%s?track=%s&lang=en", p.endpoints["GET_TRACKS"], trackID)
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return getTrackResult{}, ErrProviderFailure
	}

	resp, err := client.Do(request)

	if err != nil {
		return getTrackResult{}, ErrProviderFailure
	}

	err = p.handleError(resp)

	if err != nil {
		return getTrackResult{}, err
	}

	defer resp.Body.Close()

	var result struct {
		Track struct {
			ID     string
			Title  string
			Albums []struct {
				ID int
			}
			Artists []struct {
				ID   int
				Name string
			}
		}
	}

	var resultArtists []struct {
		ID   string
		Name string
	}

	err = json.NewDecoder(resp.Body).Decode(&result)

	for _, artist := range result.Track.Artists {
		resultArtists = append(resultArtists, struct {
			ID   string
			Name string
		}{
			ID:   strconv.Itoa(artist.ID),
			Name: artist.Name,
		})
	}

	if err != nil {
		return getTrackResult{}, ErrProviderFailure
	}

	return getTrackResult{
		ID:      result.Track.ID,
		Name:    result.Track.Title,
		Artists: resultArtists,
		Album: struct{ ID string }{
			ID: strconv.Itoa(result.Track.Albums[0].ID),
		},
	}, nil
}

func (p *yandexProvider) GetAlbum(albumID string) (getAlbumResult, error) {
	url := fmt.Sprintf("%s?album=%s&lang=en", p.endpoints["GET_ALBUMS"], albumID)
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return getAlbumResult{}, ErrProviderFailure
	}

	resp, err := client.Do(request)

	if err != nil {
		return getAlbumResult{}, ErrProviderFailure
	}

	err = p.handleError(resp)

	if err != nil {
		return getAlbumResult{}, err
	}

	defer resp.Body.Close()

	var result getYandexAlbumResult

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return getAlbumResult{}, ErrProviderFailure
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
		return getArtistResult{}, ErrProviderFailure
	}

	resp, err := client.Do(request)

	if err != nil {
		return getArtistResult{}, ErrProviderFailure
	}

	err = p.handleError(resp)

	if err != nil {
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
		return getArtistResult{}, ErrProviderFailure
	}

	return getArtistResult{
		ID:   result.Artist.ID,
		Name: result.Artist.Name,
	}, nil
}

func (p *yandexProvider) Search(name, searchType string) (map[string]string, error) {
	query := url.QueryEscape(name)
	searchURL := fmt.Sprintf("%s?text=%s&type=%s&lang=en", p.endpoints["SEARCH"], query, searchType)

	request, err := http.NewRequest("GET", searchURL, nil)

	if err != nil {
		return nil, ErrProviderFailure
	}

	resp, err := client.Do(request)

	if err != nil {
		return nil, ErrProviderFailure
	}

	err = p.handleError(resp)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result struct {
		Tracks struct {
			Items []struct {
				ID     int
				Albums []struct {
					ID int
				}
			}
		}
		Albums struct {
			Items []struct {
				ID int
			}
		}
		Artists struct {
			Items []struct {
				ID int
			}
		}
	}

	// b, _ := ioutil.ReadAll(resp.Body)
	// log.Println(string(b)) //TODO: can be readed once?

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		log.Fatal(err)
		return nil, ErrProviderFailure
	}

	trackID := ""
	albumID := ""
	artistID := ""

	log.Println("result.Artists.Items[0].ID")
	log.Println(result.Artists.Items[0].ID)

	switch searchType {
	case "track":
		trackID = strconv.Itoa(result.Tracks.Items[0].ID)
		albumID = strconv.Itoa(result.Tracks.Items[0].Albums[0].ID)
	case "album":
		albumID = strconv.Itoa(result.Albums.Items[0].ID)
	case "artist":
		artistID = strconv.Itoa(result.Artists.Items[0].ID)
	}

	return map[string]string{
		"track":  trackID,
		"album":  albumID,
		"artist": artistID,
	}, nil
}

func (p *yandexProvider) handleError(resp *http.Response) error {
	log.Print("resp.StatusCode")
	log.Println(resp.StatusCode)

	if resp.StatusCode != 200 {
		if resp.StatusCode == 400 {
			return ErrBadRequest
		} else if resp.StatusCode == 404 {
			return ErrNotFound
		} else {
			return ErrProviderFailure
		}
	}

	return nil
}

func (p *yandexProvider) GenerateLink(IDs map[string]string, linkType string) (string, error) {
	switch linkType {
	case "track":
		return fmt.Sprintf("https://music.yandex.ru/album/%s/track/%s", IDs["album"], IDs["track"]), nil
	case "album":
		return fmt.Sprintf("https://music.yandex.ru/album/%s", IDs["album"]), nil
	case "artist":
		return fmt.Sprintf("https://music.yandex.ru/artist/%s", IDs["artist"]), nil
	}

	return "", ErrWrongLinkType
}
