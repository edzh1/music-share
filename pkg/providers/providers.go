package providers

import (
	"net/http"
	"time"
)

type providerInterface interface {
	getTrack(string) string
}

var timeout = time.Duration(5 * time.Second)
var client = http.Client{
	Timeout: timeout,
}

type provider struct {
	name      string
	token     string
	endpoints map[string]string
}

func (p *provider) GetTrack(string) string {
	return "string"
}
