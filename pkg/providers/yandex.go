package providers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type yandexProvider struct {
	provider
}

//Yandex provider
var Yandex = &yandexProvider{
	provider: provider{
		Name: "yandex",
		endpoints: map[string]string{
			"GET_TRACK": "https://music.yandex.ru/handlers/track.jsx",
			"GET_ALBUM": "https://music.yandex.ru/handlers/album.jsx",
			"SEARCH":    "https://music.yandex.ru/handlers/music-search.jsx",
		},
	},
}

func (p *yandexProvider) GetTrack(trackID string) string {
	url := fmt.Sprintf("%s?track=%s&lang=en", p.endpoints["GET_TRACK"], trackID)
	request, err := http.NewRequest("GET", url, nil)

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

func (p *yandexProvider) GetAlbum(albumID string) string {
	url := fmt.Sprintf("%s?album=%s&lang=en", p.endpoints["GET_ALBUM"], albumID)
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
