package urlparser

import (
	"errors"
	"regexp"
)

//URLParser can get provider name, type of a link (track, album, etc) from an url string
type URLParser struct {
	Providers []string
}

// TODO: configurable providers
// TODO: make regexps more accurate
var regexpMap = map[string]map[string]*regexp.Regexp{
	"yandex": {
		"domain": regexp.MustCompile("music\\.yandex\\.ru"),
		"track":  regexp.MustCompile("/album/\\d+/track/\\d+"),
		"album":  regexp.MustCompile("/album/\\d+($|\\?)"),
		"artist": regexp.MustCompile("/artist/\\d+($|\\?)"),
	},
	"spotify": {
		"domain": regexp.MustCompile("open\\.spotify\\.com"),
		"track":  regexp.MustCompile("/track/\\w+"),
		"album":  regexp.MustCompile("/album/\\w+"),
		"artist": regexp.MustCompile("/artist/\\w+"),
	},
}

//GetProvider returns provider from an url string
func (p *URLParser) GetProvider(url string) (string, error) {
	for _, provider := range p.Providers {
		pattern := regexpMap[provider]["domain"]

		if pattern.MatchString(url) {
			return provider, nil
		}
	}

	return "", errors.New("URLParser: unknown provider")
}

//GetLinkType returns link type (track, album, etc) from an url string
func (p *URLParser) GetLinkType(url string) (string, error) {
	provider, err := p.GetProvider(url)

	if err != nil {
		return "", err
	}

	pattern := regexpMap[provider]["track"]

	if pattern.MatchString(url) {
		return "track", nil
	}

	pattern = regexpMap[provider]["album"]

	if pattern.MatchString(url) {
		return "album", nil
	}

	pattern = regexpMap[provider]["artist"]

	if pattern.MatchString(url) {
		return "artist", nil
	}

	return "", errors.New("URLParser: unknown link type")
}
