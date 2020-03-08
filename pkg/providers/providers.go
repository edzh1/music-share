package providers

import (
	"net/http"
	"time"
)

type providerInterface interface {
	GetTrack(string) (string, error)
	GetAlbum(string) (string, error)
	GetArtist(string) (string, error)
}

var timeout = time.Duration(5 * time.Second)
var client = http.Client{
	Timeout: timeout,
}

//Provider struct
type Provider struct {
	Name              string
	ClientToken       string
	apiToken          string
	endpoints         map[string]string
	providerInterface //make this work
}
