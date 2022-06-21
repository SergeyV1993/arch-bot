package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
)

const collectionName = "locations"

const defaultLimit = 25

type MongoRepository struct {
	db *mongo.Database
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{db: db}
}

func (mr *MongoRepository) GetBuildsByDistance(ctx context.Context, longitude, latitude string, distance int) ([]ArchitectBuilds, error) {
	abs := make([]ArchitectBuilds, 0)

	lat, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		return nil, err
	}
	long, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		return nil, err
	}
	location := NewLocation(long, lat)

	collections := mr.db.Collection(collectionName)
	pipeline, err := mr.createPipeline(location, distance)
	if err != nil {
		return nil, err
	}
	res, err := collections.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer res.Close(ctx)

	var ab ArchitectBuilds
	for res.Next(ctx) {
		err = res.Decode(&ab)
		if err != nil {
			return nil, err
		}

		abs = append(abs, ab)
	}

	return abs, nil
}

func (mr *MongoRepository) createPipeline(location Location, distance int) (mongo.Pipeline, error) {
	geoNearStage := bson.D{
		{"$geoNear",
			bson.D{
				{"includeLocs", "location"},
				{"distanceField", "distance"},
				{"maxDistance", distance},
				{"maxDistance", true},
				{"near", location},
			}},
	}

	limitStage := bson.D{{"$limit", defaultLimit}}
	sortStage := bson.D{{"$sort", bson.D{{"distance", 1}}}}

	var pipeline mongo.Pipeline
	pipeline = append(pipeline, geoNearStage, limitStage, sortStage)

	return pipeline, nil
}
