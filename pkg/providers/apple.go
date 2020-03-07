package providers

type appleProvider struct {
	provider
}

//Apple provider
var Apple = &appleProvider{
	provider: provider{
		name: "apple",
	},
}

// func (p *appleProvider) getTrack(string) string {
// 	return "1"
// }
