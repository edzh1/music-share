package mongo

import (
	"context"
	"log"

	"github.com/edzh1/music-share/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//AlbumModel for a single album
type AlbumModel struct {
	Client *mongo.Client
}

//Insert single album into db
func (m *AlbumModel) Insert(name, spotifyID string) (int, error) {
	collection := m.Client.Database("music-share").Collection("music")
	res, err := collection.InsertOne(context.Background(), bson.M{"name": name, "spotifyID": spotifyID})

	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	log.Println(res)

	return 1, nil
}

//Get single album from db
func (m *AlbumModel) Get(filter bson.M) (models.Album, error) {
	var result models.Album

	collection := m.Client.Database("music-share").Collection("music")
	err := collection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		log.Fatal(err)
		return models.Album{}, err
	}

	log.Println(result)

	return result, nil
}

//Delete single album from db
func (m *AlbumModel) Delete(spotifyID string) (int, error) {
	collection := m.Client.Database("music-share").Collection("music")
	result, err := collection.DeleteOne(context.Background(), bson.M{"spotifyID": spotifyID})

	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	log.Println(result)

	return 1, nil
}
