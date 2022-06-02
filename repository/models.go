package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	ID        primitive.ObjectID `bson:"id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Msg       string             `bson:"msg"`
}

type ArchitectBuilds struct {
	Name           string `bson:"name"`
	Address        string `bson:"address"`
	LinkMapAddress string `bson:"link_map_address"`

	Link        string `bson:"link"`
	Description string `bson:"description"`

	Location Location `bson:"location"`
	Distance float64  `bson:"distance"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type Location struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}

func NewLocation(long, lat float64) Location {
	return Location{
		"Point",
		[]float64{long, lat},
	}
}
