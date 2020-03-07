package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//TrackModel for a single track
type TrackModel struct {
	Client *mongo.Client
	Ctx    *context.Context
}

//Insert single track into db
func (m *TrackModel) Insert(name, spotifyID string) (int, error) {
	collection := m.Client.Database("music-share").Collection("tracks")
	res, err := collection.InsertOne(*m.Ctx, bson.M{"name": name, "spotifyID": spotifyID})

	if err != nil {
		log.Fatal(err)
	}

	log.Println(res)

	return 1, nil
}

//Get single track from db
func (m *TrackModel) Get(spotifyID string) (int, error) {
	var result struct {
		Name string
	}

	collection := m.Client.Database("music-share").Collection("tracks")
	err := collection.FindOne(*m.Ctx, bson.M{"spotifyID": spotifyID}).Decode(&result)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(result)

	return 1, nil
}

//Delete single track from db
func (m *TrackModel) Delete(spotifyID string) (int, error) {
	collection := m.Client.Database("music-share").Collection("tracks")
	result, err := collection.DeleteOne(*m.Ctx, bson.M{"spotifyID": spotifyID})

	if err != nil {
		log.Fatal(err)
	}

	log.Println(result)

	return 1, nil
}
