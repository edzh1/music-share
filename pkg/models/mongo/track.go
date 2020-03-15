package mongo

import (
	"context"
	"log"

	"github.com/edzh1/music-share/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//TrackModel for a single track
type TrackModel struct {
	Client *mongo.Client
}

//Insert single track into db
func (m *TrackModel) Insert(document bson.M) (int, error) {
	collection := m.Client.Database("music-share").Collection("tracks")
	res, err := collection.InsertOne(context.Background(), document)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(res)

	return 1, nil
}

//Get single track from db
func (m *TrackModel) Get(filter bson.M) (models.Track, error) {
	var result models.Track

	collection := m.Client.Database("music-share").Collection("tracks")
	err := collection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		return models.Track{}, err
	}

	return result, nil
}

//Delete single track from db
func (m *TrackModel) Delete(spotifyID string) (int, error) {
	collection := m.Client.Database("music-share").Collection("tracks")
	result, err := collection.DeleteOne(context.Background(), bson.M{"spotifyID": spotifyID})

	if err != nil {
		log.Fatal(err)
	}

	log.Println(result)

	return 1, nil
}
