package mongo

import (
	"context"
	"log"

	"github.com/edzh1/music-share/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//ArtistModel for a single artist
type ArtistModel struct {
	Client *mongo.Client
}

var artistsCollection = "artists"

//Insert single artist into db
func (m *ArtistModel) Insert(document bson.M) (int, error) {
	collection := m.Client.Database("music-share").Collection(artistsCollection)
	res, err := collection.InsertOne(context.Background(), document)

	if err != nil {
		return 0, err
	}

	log.Println(res)

	return 1, nil
}

//Get single artist from db
func (m *ArtistModel) Get(filter bson.M) (models.Artist, error) {
	var result models.Artist

	collection := m.Client.Database("music-share").Collection(artistsCollection)
	err := collection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		return models.Artist{}, err
	}

	log.Println(result)

	return result, nil
}

//Delete single artist from db
func (m *ArtistModel) Delete(spotifyID string) (int, error) {
	collection := m.Client.Database("music-share").Collection(artistsCollection)
	result, err := collection.DeleteOne(context.Background(), bson.M{"spotifyID": spotifyID})

	if err != nil {
		return 0, err
	}

	log.Println(result)

	return 1, nil
}
