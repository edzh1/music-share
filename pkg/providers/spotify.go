package providers

type spotifyProvider struct {
	provider
}

var Spotify = &spotifyProvider{
	provider: provider{
		name: "spotify",
	},
}

func (p *spotifyProvider) getTrack(string) string {
	return "1"
}
