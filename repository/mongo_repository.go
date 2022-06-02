package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"strconv"
	"time"
)

const collectionName = "locations"

type MongoRepository struct {
	db *mongo.Database
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{db: db}
}

//TODO debug
func (mr *MongoRepository) SetBuilds(ctx context.Context, build ArchitectBuilds) error {
	msg := ArchitectBuilds{
		Name:           "Особняк Миндовского",
		Address:        "Поварская, 44/2с",
		LinkMapAddress: "https://yandex.ru/maps?mode=search&text=55.756828,37.588665",
		Link:           "https://yandex.ru/",
		Description:    "Особняк представительства Новой Зеландии. Роскошный модерн ХХ века",
		Location:       NewLocation(37.588665, 55.756828),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := mr.createIndex(ctx)
	if err != nil {
		return err
	}

	collections := mr.db.Collection(collectionName)
	_, err = collections.InsertOne(ctx, msg)
	if err != nil {
		return err
	}

	return nil
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
	/*
		filter := bson.D{
			{"location",
				bson.D{
					{"$near", bson.D{
						{"$geometry", location},
						{"$maxDistance", distance},
					}},
				}},
		}

		res, err := collections.Find(ctx, filter)
		if err != nil {
			return nil, err
		}

		//defer res.Close(ctx)

		var ab ArchitectBuilds
		for res.Next(ctx) {
			err = res.Decode(&ab)
			if err != nil {
				return nil, err
			}

			abs = append(abs, ab)
		}
	*/
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

func (mr *MongoRepository) createIndex(ctx context.Context) error {
	indexOpts := options.CreateIndexes().SetMaxTime(time.Second * 10)

	pointIndexModel := mongo.IndexModel{
		Options: options.Index(),
		Keys:    bsonx.MDoc{"location": bsonx.String("2dsphere")},
	}

	_, err := mr.db.Collection(collectionName).Indexes().CreateOne(ctx, pointIndexModel, indexOpts)
	if err != nil {
		return err
	}

	return nil
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

	var pipeline mongo.Pipeline
	pipeline = append(pipeline, geoNearStage)

	return pipeline, nil
}
