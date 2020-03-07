package main

import (
	"github.com/edzh1/music-share/providers"
)

func main() {
	message := providers.Spotify.getTrack("string")

	println(message)
}
