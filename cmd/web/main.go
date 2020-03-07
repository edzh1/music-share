package main

import (
	"flag"

	"github.com/edzh1/music-share/pkg/providers"
)

func main() {
	spotifyCredntials := flag.String("spotifyCredntials", "", "Base64 encoded client_id:clent_secret")
	flag.Parse()

	providers.Spotify.ClientToken = *spotifyCredntials
	message := providers.Spotify.GetTrack("7rqWfVnNo2hyCpSpCpEYFj")
	// message := providers.Spotify.Auth()

	println(message)
}
