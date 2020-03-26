package main

import (
	"context"
	"flag"
	"log"
	"time"

	mongoModels "github.com/edzh1/music-share/pkg/models/mongo"
	"github.com/edzh1/music-share/pkg/providers"
	"github.com/edzh1/music-share/pkg/urlparser"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type providerMap map[string]providers.ProviderInterface

type application struct {
	tracks    *mongoModels.TrackModel
	albums    *mongoModels.AlbumModel
	artists   *mongoModels.ArtistModel
	providers providerMap
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	defer cancel()

	var app = &application{
		tracks:  &mongoModels.TrackModel{Client: client},
		albums:  &mongoModels.AlbumModel{Client: client},
		artists: &mongoModels.ArtistModel{Client: client},
		providers: providerMap{
			"spotify": providers.Spotify,
			"yandex":  providers.Yandex,
		},
	}

	spotifyCredentials := flag.String("spotifyCredentials", "", "Base64 encoded client_id:clent_secret")
	flag.Parse()

	providers.Spotify.ClientToken = *spotifyCredentials

	providerParser := &urlparser.URLParser{
		Providers: []string{"yandex", "spotify"},
	}

	provider, err := providerParser.GetProvider("https://open.spotify.com/track/4DZBk2qmeAvZtSmhSayaXh")
	// provider, err := providerParser.GetProvider("https://music.yandex.ru/album/10073206/track/63360735")

	if err != nil {
		log.Fatal(err)
	}

	// _, _ = app.getTrack("28358063", app.providers[provider])
	// _, _ = app.getAlbum("9899790", app.providers[provider])
	// _, _ = app.getArtist("1768379", app.providers[provider])
	// _, _ = app.getTrack("4DZBk2qmeAvZtSmhSayaXh", app.providers[provider])
	_, _ = providers.Yandex.GetEntityID("https://music.yandex.ru/album/10073206/track/63360735", "track")
	// _, _ = app.getAlbum("7CpFDGlrIfHLBmCBPvTSLU?si=GHCzKcKgTeKlRgzXDRdw8w", app.providers[provider])
	_, _ = app.getArtist("24eQxPRLv3UMwEIo6mawVW", app.providers[provider])

	linkType, err := providerParser.GetLinkType("https://open.spotify.com/track/7rqWfVnNo2hyCpSpCpEYFj?si=GHCzKcKgTeKlRgzXDRdw8w")
	// linkType, err := providerParser.GetLinkType("https://music.yandex.ru/album/10073206/track/63360735")

	if err != nil {
		log.Fatal(err)
	}

	println(provider)
	println(linkType)
}
