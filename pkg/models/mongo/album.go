package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//AlbumModel for a single album
type AlbumModel struct {
	Client *mongo.Client
}

//Insert single track into db
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

//Get single track from db
func (m *AlbumModel) Get(spotifyID string) (string, error) {
	var result struct {
		Name string
	}

	collection := m.Client.Database("music-share").Collection("music")
	err := collection.FindOne(context.Background(), bson.M{"spotifyID": spotifyID}).Decode(&result)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	log.Println(result)

	return result.Name, nil
}

//Delete single track from db
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
