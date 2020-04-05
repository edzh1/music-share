package main

import (
	"encoding/json"
	"net/http"

	"github.com/edzh1/music-share/pkg/providers"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func (app *application) handleHTTPError(w http.ResponseWriter, err error) {
	if err == providers.ErrBadRequest {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	} else if err == providers.ErrNotFound {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	app.errorLog.Fatal(err)
}

func (app *application) handleLink(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
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

	result := make(map[string]string)

	for p := range app.providers {
		providerIDs := make(map[string]string)

		switch linkType {
		case "track":
			res, err := app.getTrack(ID, provider)

			switch p {
			case "spotify":
				providerIDs["track"] = res.SpotifyID
			case "yandex":
				providerIDs["track"] = res.YandexID
				providerIDs["album"] = res.YandexAlbumID
			}

			if err != nil {
				app.handleHTTPError(w, err)
				return
			}
		case "album":
			res, err := app.getAlbum(ID, provider)

			switch p {
			case "spotify":
				providerIDs["album"] = res.SpotifyID
			case "yandex":
				providerIDs["album"] = res.YandexID
			}

			if err != nil {
				app.handleHTTPError(w, err)
				return
			}
		case "artist":
			res, err := app.getArtist(ID, provider)

			switch p {
			case "spotify":
				providerIDs["artist"] = res.SpotifyID
			case "yandex":
				providerIDs["artist"] = res.YandexID
			}

			if err != nil {
				app.handleHTTPError(w, err)
				return
			}
		}

		result[p], err = app.providers[p].GenerateLink(providerIDs, linkType)

		if err != nil {
			app.handleHTTPError(w, err)
			return
		}
	}

	b, err := json.Marshal(result)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		app.errorLog.Fatal(err)
	}

	w.Write(b)

	return
}
