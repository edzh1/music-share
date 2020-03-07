package providers

type providerInterface interface {
	getTrack(string) string
}

type provider struct {
	name string
}

func (p *provider) getTrack(string) string {
	return "string"
}
