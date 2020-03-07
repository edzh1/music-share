package providers

type yandexProvider struct {
	provider
}

//Yandex provider
var Yandex = &yandexProvider{
	provider: provider{
		Name: "yandex",
	},
}
