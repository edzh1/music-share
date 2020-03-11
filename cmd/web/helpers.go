package main

import (
	// mongoModels "github.com/edzh1/music-share/pkg/models/mongo"
	"encoding/json"

	"github.com/edzh1/music-share/pkg/models"
	// "github.com/edzh1/music-share/pkg/urlparser"

	"log"

	"github.com/edzh1/music-share/pkg/providers"
)

func (app *application) getTrack(ID string, provider providers.ProviderInterface) (models.Track, error) {
	// result, err := app.tracks.Get(ID)

	// if err != nil && err != mongo.ErrNoDocuments {
	// 	log.Fatal(err)
	// 	return models.Track{}, err
	// }

	// if err == nil {
	// 	return result, nil
	// }

	providerResult, err := provider.GetTrack(ID)
	out, err := json.Marshal(providerResult)

	log.Println(string(out))

	if err != nil {
		log.Fatal(err)
		return models.Track{}, err
	}

	return models.Track{}, nil
}

func (app *application) getAlbum(ID string, provider providers.ProviderInterface) (models.Album, error) {
	// result, err := app.albums.Get(ID)

	// if err != nil && err != mongo.ErrNoDocuments {
	// 	log.Fatal(err)
	// 	return "", err
	// }

	// if err == nil {
	// 	return result, nil
	// }

	providerResult, err := provider.GetAlbum(ID)
	out, err := json.Marshal(providerResult)

	log.Println(string(out))

	if err != nil {
		log.Fatal(err)
		return models.Album{}, err
	}

	return models.Album{}, nil
}

func (app *application) getArtist(ID string, provider providers.ProviderInterface) (models.Artist, error) {
	// result, err := app.tracks.Get(ID)

	// if err != nil && err != mongo.ErrNoDocuments {
	// 	log.Fatal(err)
	// 	return models.Track{}, err
	// }

	// if err == nil {
	// 	return result, nil
	// }

	providerResult, err := provider.GetArtist(ID)
	out, err := json.Marshal(providerResult)

	log.Println(string(out))

	if err != nil {
		log.Fatal(err)
		return models.Artist{}, err
	}

	return models.Artist{}, nil
}
