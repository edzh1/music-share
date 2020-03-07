package providers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type spotifyProvider struct {
	provider
}

//Spotify provider
var Spotify = &spotifyProvider{
	provider: provider{
		name:  "spotify",
		token: "",
		endpoints: map[string]string{
			"GET_TRACK": "https://api.spotify.com/v1/tracks",
			"GET_ALBUM": "https://api.spotify.com/v1/albums",
			"SEARCH":    "https://api.spotify.com/v1/search",
			"AUTH":      "https://accounts.spotify.com/api/token",
		},
	},
}

func (p *spotifyProvider) GetTrack(trackID string) string {
	url := fmt.Sprintf("%s/%s", Spotify.endpoints["GET_TRACK"], trackID)
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("Authorization", p.token)

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
