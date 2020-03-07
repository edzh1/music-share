package providers

type appleProvider struct {
	provider
}

//Apple provider
var Apple = &appleProvider{
	provider: provider{
		Name: "apple",
	},
}

// func (p *appleProvider) getTrack(string) string {
// 	return "1"
// }
