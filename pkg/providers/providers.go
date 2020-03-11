package providers

import (
	"net/http"
	"time"
)

type ProviderInterface interface {
	GetTrack(string) (getTrackResult, error)
	GetAlbum(string) (getAlbumResult, error)
	GetArtist(string) (getArtistResult, error)
}

var timeout = time.Duration(5 * time.Second)
var client = http.Client{
	Timeout: timeout,
}

type artist struct {
	ID   string
	Name string
}

type getTrackResult struct {
	ID      string
	Name    string
	Artists []*struct {
		ID   string
		Name string
	}
}

type getAlbumResult struct {
	ID   string
	Name string
}

type getArtistResult struct {
	ID   string
	Name string
}

//Provider struct
type Provider struct {
	Name              string
	ClientToken       string
	apiToken          string
	endpoints         map[string]string
	ProviderInterface //make this work
}
