package main

import (
	// mongoModels "github.com/edzh1/music-share/pkg/models/mongo"
	// "github.com/edzh1/music-share/pkg/urlparser"

	"log"

	"github.com/edzh1/music-share/pkg/providers"
	"go.mongodb.org/mongo-driver/mongo"
)

func (app *application) getTrack(ID string, provider *providers.Provider) (string, error) {
	result, err := app.track.Get(ID)

	if err != nil && err != mongo.ErrNoDocuments {
		log.Fatal(err)
		return "", err
	}

	if err == nil {
		return result, nil
	}

	result, err = provider.GetTrack(ID)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return result, nil
}
