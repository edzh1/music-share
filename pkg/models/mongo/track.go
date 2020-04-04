package mongo

import (
	"context"

	"github.com/edzh1/music-share/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//TrackModel for a single track
type TrackModel struct {
	Client *mongo.Client
}

var tracksCollection = "tracks"

//Insert single track into db
func (m *TrackModel) Insert(document bson.M) (int, error) {
	collection := m.Client.Database("music-share").Collection(tracksCollection)
	_, err := collection.InsertOne(context.Background(), document)

	if err != nil {
		return 0, err
	}

	return 1, nil
}

//Get single track from db
func (m *TrackModel) Get(filter bson.M) (models.Track, error) {
	var result models.Track

	collection := m.Client.Database("music-share").Collection(tracksCollection)
	err := collection.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		return models.Track{}, err
	}

	return result, nil
}

//Delete single track from db
func (m *TrackModel) Delete(spotifyID string) (int, error) {
	collection := m.Client.Database("music-share").Collection(tracksCollection)
	_, err := collection.DeleteOne(context.Background(), bson.M{"spotifyID": spotifyID})

	if err != nil {
		return 0, err
	}

	return 1, nil
}
