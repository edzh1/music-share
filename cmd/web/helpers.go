package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/edzh1/music-share/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	// "github.com/edzh1/music-share/pkg/urlparser"

	"log"

	"github.com/edzh1/music-share/pkg/providers"
)

func (app *application) handleLink(w http.ResponseWriter, r *http.Request) {
	URL := r.URL.Query().Get("url")
	providerName, err := app.providerParser.GetProvider(URL)

	log.Println(URL)

	if err != nil {
		return
	}

	linkType, err := app.providerParser.GetLinkType(URL)

	if err != nil {
		return
	}

	provider := app.providers[providerName]

	ID, err := provider.GetEntityID(URL, linkType)

	if err != nil {
		return
	}

	switch linkType {
	case "track":
		res, err := app.getTrack(ID, provider)

		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println(res.SpotifyID)

		b, err := json.Marshal(res)

		if err != nil {
			return
		}

		w.Write(b)
	case "album":
		app.getAlbum(ID, provider)
	case "artist":
		app.getArtist(ID, provider)
	}

	return
}

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

	log.Println("before EOF")

	providerResult, err := provider.GetTrack(ID)

	log.Println("before EOF2")
	log.Println(err)

	if err != nil {
		return models.Track{}, err
	}

	log.Println("PIZDA!")

	out, err := json.Marshal(providerResult)

	log.Println("PIZDA")
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

		var result models.Track

		log.Println(newTrack)

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
