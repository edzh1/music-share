package providers

import (
	"errors"
	"net/http"
	"time"
)

var ErrBadRequest = errors.New("providers: 400")
var ErrNotFound = errors.New("providers: 404")
var ErrProviderFailure = errors.New("providers: 500")
var ErrWrongSearchType = errors.New("wrong search type")
var ErrAuth = errors.New("auth error")

type ProviderInterface interface {
	GetName() string
	GetTrack(string) (getTrackResult, error)
	GetAlbum(string) (getAlbumResult, error)
	GetArtist(string) (getArtistResult, error)
	Search(name, searchType string) (string, error)
	GetEntityID(url, entity string) (string, error)
}

var timeout = time.Duration(15 * time.Second)
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
	ID      string
	Name    string
	Artists []*struct {
		ID   string
		Name string
	}
}

type getArtistResult struct {
	ID   string
	Name string
}

//Provider struct
type Provider struct {
	Name      string
	endpoints map[string]string
}
