package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/edzh1/music-share/pkg/models"
	"github.com/edzh1/music-share/pkg/providers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (app *application) handleLink(w http.ResponseWriter, r *http.Request) {
	URL := r.URL.Query().Get("url")
	providerName, err := app.providerParser.GetProvider(URL)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	linkType, err := app.providerParser.GetLinkType(URL)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	provider := app.providers[providerName]

	ID, err := provider.GetEntityID(URL, linkType)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var res interface{}

	switch linkType {
	case "track":
		res, err = app.getTrack(ID, provider)
	case "album":
		res, err = app.getAlbum(ID, provider)
	case "artist":
		res, err = app.getArtist(ID, provider)
	}

	if err != nil {
		if err == providers.ErrBadRequest {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		} else if err == providers.ErrProviderFailure {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		log.Fatal(err)
	}

	b, err := json.Marshal(res)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	w.Write(b)

	return
}

func (app *application) getTrack(ID string, provider providers.ProviderInterface) (models.Track, error) {
	filter := bson.M{fmt.Sprintf("%sID", provider.GetName()): ID}

	result, err := app.tracks.Get(filter)

	if err != nil && err != mongo.ErrNoDocuments {
		return models.Track{}, err
	}

	if err == nil {
		return result, nil
	}

	providerResult, err := provider.GetTrack(ID)

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
				providerID, err := providerValue.Search(fmt.Sprintf("%s %s", providerResult.Name, artists), "track")

				if err != nil {
					providerID = ""
				}

				newTrack[fmt.Sprintf("%sID", providerKey)] = providerID
			}
		}

		_, _ = app.tracks.Insert(newTrack)

		var result models.Track

		bsonBytes, _ := bson.Marshal(newTrack)
		bson.Unmarshal(bsonBytes, &result)

		return result, nil
	}

	return models.Track{}, nil
}

func (app *application) getAlbum(ID string, provider providers.ProviderInterface) (models.Album, error) {
	filter := bson.M{fmt.Sprintf("%sID", provider.GetName()): ID}

	result, err := app.albums.Get(filter)

	if err != nil && err != mongo.ErrNoDocuments {
		return models.Album{}, err
	}

	if err == nil {
		return result, nil
	}

	providerResult, err := provider.GetAlbum(ID)

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
				providerID, err := providerValue.Search(fmt.Sprintf("%s - %s", providerResult.Name, artists), "album")

				if err != nil {
					providerID = ""
				}

				newAlbum[fmt.Sprintf("%sID", providerKey)] = providerID
			}
		}

		_, _ = app.albums.Insert(newAlbum)

		var result models.Album

		bsonBytes, _ := bson.Marshal(newAlbum)
		bson.Unmarshal(bsonBytes, &result)

		return result, nil
	}

	return models.Album{}, nil
}

func (app *application) getArtist(ID string, provider providers.ProviderInterface) (models.Artist, error) {
	filter := bson.M{fmt.Sprintf("%sID", provider.GetName()): ID}

	result, err := app.artists.Get(filter)

	if err != nil && err != mongo.ErrNoDocuments {
		return models.Artist{}, err
	}

	if err == nil {
		return result, nil
	}

	providerResult, err := provider.GetArtist(ID)

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
				providerID, err := providerValue.Search(providerResult.Name, "artist")

				if err != nil {
					providerID = ""
				}

				newArtist[fmt.Sprintf("%sID", providerKey)] = providerID
			}
		}

		_, _ = app.artists.Insert(newArtist)

		var result models.Artist

		bsonBytes, _ := bson.Marshal(newArtist)
		bson.Unmarshal(bsonBytes, &result)

		return result, nil
	}

	return models.Artist{}, nil
}
