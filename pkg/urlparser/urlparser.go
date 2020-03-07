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
var regexpMap = map[string]*regexp.Regexp{
	"yandex":  regexp.MustCompile("music\\.yandex\\.ru"),
	"spotify": regexp.MustCompile("open\\.spotify\\.com"),
}

//GetProvider returns provider from an url string
func (p *URLParser) GetProvider(url string) (string, error) {
	for _, provider := range p.Providers {
		pattern := regexpMap[provider]

		if pattern.MatchString(url) {
			return provider, nil
		}
	}

	return "", errors.New("unknown provider")
}
