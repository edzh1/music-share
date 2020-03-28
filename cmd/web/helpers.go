package main

import (
	"fmt"
	"strings"

	"github.com/edzh1/music-share/pkg/models"
	"github.com/edzh1/music-share/pkg/providers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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

	var providerIDs map[string]string

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

		if provider.GetName() == "yandex" {
			newTrack["YandexAlbumID"] = providerResult.Album.ID
		}

		for providerKey, providerValue := range app.providers {
			if providerKey != provider.GetName() {
				providerIDs, err = providerValue.Search(fmt.Sprintf("%s %s", providerResult.Name, artists), "track")

				if err != nil {
					providerIDs = nil
				}

				newTrack[fmt.Sprintf("%sID", providerKey)] = providerIDs["track"]

				if providerKey == "yandex" {
					newTrack["YandexAlbumID"] = providerIDs["album"]
				}
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
				providerIDs, err := providerValue.Search(fmt.Sprintf("%s - %s", providerResult.Name, artists), "album")

				if err != nil {
					providerIDs = nil
				}

				newAlbum[fmt.Sprintf("%sID", providerKey)] = providerIDs["album"]
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
				providerIDs, err := providerValue.Search(providerResult.Name, "artist")

				if err != nil {
					providerIDs = nil
				}

				newArtist[fmt.Sprintf("%sID", providerKey)] = providerIDs["artist"]
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
