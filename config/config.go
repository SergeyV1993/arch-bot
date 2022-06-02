package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

type Env struct {
	TelegramToken string

	DbUrl               string
	DbName              string
	DbConnectionTimeout time.Duration

	RadiusMeters int
}

func NewEnv(token, dbUrl, dbName string, dbConnectionTimeout time.Duration, radius int) *Env {
	return &Env{
		TelegramToken:       token,
		DbUrl:               dbUrl,
		DbName:              dbName,
		DbConnectionTimeout: dbConnectionTimeout,
		RadiusMeters:        radius,
	}
}

func InitEnvs() (*Env, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		return nil, err
	}

	dbCon, err := strconv.Atoi(os.Getenv("DB_CONNECTION_TIMEOUT"))
	if err != nil {
		return nil, err
	}
	dbConTimeout := time.Duration(dbCon) * time.Second

	radius, err := strconv.Atoi(os.Getenv("RADIUS"))
	if err != nil {
		return nil, err
	}

	e := NewEnv(
		os.Getenv("TELEGRAMM_TOKEN"),
		os.Getenv("DB_URL"),
		os.Getenv("DB_NAME"),
		dbConTimeout,
		radius,
	)

	return e, nil
}
