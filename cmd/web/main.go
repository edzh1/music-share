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

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	var track = &mongoModels.TrackModel{Client: client, Ctx: &ctx}
	// track.Insert("Technoid", "7CpFDGlrIfHLBmCBPvTSLU")
	// track.Get("7CpFDGlrIfHLBmCBPvTSLU")
	// track.Delete("7CpFDGlrIfHLBmCBPvTSLU")

	spotifyCredntials := flag.String("spotifyCredntials", "", "Base64 encoded client_id:clent_secret")
	flag.Parse()

	providerParser := &urlparser.URLParser{
		Providers: []string{"yandex", "spotify"},
	}

	provider, err := providerParser.GetProvider("https://open.spotify.com/track/7rqWfVnNo2hyCpSpCpEYFj?si=GHCzKcKgTeKlRgzXDRdw8w")

	if err != nil {
		log.Fatal(err)
	}

	linkType, err := providerParser.GetLinkType("https://open.spotify.com/track/7rqWfVnNo2hyCpSpCpEYFj?si=GHCzKcKgTeKlRgzXDRdw8w")

	if err != nil {
		log.Fatal(err)
	}

	providers.Spotify.ClientToken = *spotifyCredntials
	// message := providers.Spotify.GetTrack("7rqWfVnNo2hyCpSpCpEYFj")
	// message := providers.Spotify.GetAlbum("7CpFDGlrIfHLBmCBPvTSLU")
	// message := providers.Yandex.GetTrack("62650884%3A9899790")
	// message := providers.Spotify.GetAlbum("7CpFDGlrIfHLBmCBPvTSLU")
	// message := providers.Spotify.Auth()

	println(provider)
	println(linkType)
}
