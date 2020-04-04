package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	mongoModels "github.com/edzh1/music-share/pkg/models/mongo"
	"github.com/edzh1/music-share/pkg/providers"
	"github.com/edzh1/music-share/pkg/urlparser"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type providerMap map[string]providers.ProviderInterface

type application struct {
	tracks         *mongoModels.TrackModel
	albums         *mongoModels.AlbumModel
	artists        *mongoModels.ArtistModel
	providers      providerMap
	providerParser *urlparser.URLParser
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://db:27017"))

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		log.Fatal(err)
	}

	defer cancel()

	spotifyCredentials := os.Getenv("SPOTIFY_CREDENTIALS")
	providers.Spotify.ClientToken = spotifyCredentials

	providerParser := &urlparser.URLParser{
		Providers: []string{"yandex", "spotify"},
	}

	var app = &application{
		tracks:         &mongoModels.TrackModel{Client: client},
		albums:         &mongoModels.AlbumModel{Client: client},
		artists:        &mongoModels.ArtistModel{Client: client},
		providerParser: providerParser,
		providers: providerMap{
			"spotify": providers.Spotify,
			"yandex":  providers.Yandex,
		},
	}

	// tlsConfig := &tls.Config{
	// 	PreferServerCipherSuites: true,
	// 	CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	// }

	srv := &http.Server{
		Addr:    ":4000",
		Handler: app.routes(),
		// TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// err = srv.ListenAndServe("./tls/cert.pem", "./tls/key.pem")
	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
