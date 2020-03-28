package models

//Track type for a single track
type Track struct {
	ID            string
	SpotifyID     string
	YandexID      string
	YandexAlbumID string
	Name          string
}

//Album type for a single album
type Album struct {
	ID        string
	SpotifyID string
	YandexID  string
	Name      string
}

//Artist type for a single artist
type Artist struct {
	ID        string
	SpotifyID string
	YandexID  string
	Name      string
}
