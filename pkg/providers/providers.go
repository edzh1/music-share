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
	Name      string
	endpoints map[string]string
}
