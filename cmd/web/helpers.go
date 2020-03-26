package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/edzh1/music-share/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	// "github.com/edzh1/music-share/pkg/urlparser"

	"log"

	"github.com/edzh1/music-share/pkg/providers"
)

func (app *application) getTrack(ID string, provider providers.ProviderInterface) (models.Track, error) {
	filter := bson.M{fmt.Sprintf("%sID", provider.GetName()): ID}

	result, err := app.tracks.Get(filter)

	if err != nil && err != mongo.ErrNoDocuments {
		log.Fatal(err)
		return models.Track{}, err
	}

	if err == nil {
		return result, nil
	}

	providerResult, err := provider.GetTrack(ID)

	if err != nil {
		return models.Track{}, err
	}

	out, err := json.Marshal(providerResult)

	log.Println(string(out))

	if err != nil {
		return models.Track{}, err
	}

	if result == (models.Track{}) {
		var artistSlice []string

		for _, artist := range providerResult.Artists {
			artistSlice = append(artistSlice, artist.Name)
		}

		artists := strings.Join(artistSlice, ",")

		newTrack := bson.M{
			"name":                                  providerResult.Name,
			fmt.Sprintf("%sID", provider.GetName()): ID,
		}

		for providerKey, providerValue := range app.providers {
			if providerKey != provider.GetName() {
				providerID, err := providerValue.Search(fmt.Sprintf("%s - %s", providerResult.Name, artists), "track")

				if err != nil {
					log.Fatal(err)
				}

				newTrack[fmt.Sprintf("%sID", providerKey)] = providerID
			}
		}

		_, _ = app.tracks.Insert(newTrack)
	}

	return models.Track{}, nil
}

func (app *application) getAlbum(ID string, provider providers.ProviderInterface) (models.Album, error) {
	filter := bson.M{fmt.Sprintf("%sID", provider.GetName()): ID}

	result, err := app.albums.Get(filter)

	if err != nil && err != mongo.ErrNoDocuments {
		log.Fatal(err)
		return models.Album{}, err
	}

	if err == nil {
		return result, nil
	}

	providerResult, err := provider.GetAlbum(ID)

	if err != nil {
		return models.Album{}, err
	}

	out, err := json.Marshal(providerResult)

	log.Println(string(out))

	if err != nil {
		return models.Album{}, err
	}

	if result == (models.Album{}) {
		var artistSlice []string

		for _, artist := range providerResult.Artists {
			artistSlice = append(artistSlice, artist.Name)
		}

		artists := strings.Join(artistSlice, ",")

		newAlbum := bson.M{
			"name":                                  providerResult.Name,
			fmt.Sprintf("%sID", provider.GetName()): ID,
		}

		for providerKey, providerValue := range app.providers {
			if providerKey != provider.GetName() {
				providerID, err := providerValue.Search(fmt.Sprintf("%s - %s", providerResult.Name, artists), "track")

				if err != nil {
					log.Fatal(err)
				}

				newAlbum[fmt.Sprintf("%sID", providerKey)] = providerID
			}
		}

		_, _ = app.albums.Insert(newAlbum)
	}

	return models.Album{}, nil
}

func (app *application) getArtist(ID string, provider providers.ProviderInterface) (models.Artist, error) {
	filter := bson.M{fmt.Sprintf("%sID", provider.GetName()): ID}

	result, err := app.artists.Get(filter)

	if err != nil && err != mongo.ErrNoDocuments {
		log.Fatal(err)
		return models.Artist{}, err
	}

	if err == nil {
		return result, nil
	}

	providerResult, err := provider.GetArtist(ID)

	if err != nil {
		return models.Artist{}, err
	}

	out, err := json.Marshal(providerResult)

	log.Println(string(out))

	if err != nil {
		return models.Artist{}, err
	}

	if result == (models.Artist{}) {
		newArtist := bson.M{
			"name":                                  providerResult.Name,
			fmt.Sprintf("%sID", provider.GetName()): ID,
		}

		for providerKey, providerValue := range app.providers {
			if providerKey != provider.GetName() {
				providerID, err := providerValue.Search(providerResult.Name, "artists")

				if err != nil {
					log.Fatal(err)
				}

				newArtist[fmt.Sprintf("%sID", providerKey)] = providerID
			}
		}

		_, _ = app.artists.Insert(newArtist)
	}

	return models.Artist{}, nil
}
