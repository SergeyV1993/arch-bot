package main

import (
	"ArchitectureBot/client"
	"ArchitectureBot/config"
	"ArchitectureBot/consumer"
	"ArchitectureBot/event"
	"ArchitectureBot/repository"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main() {
	ctx := context.Background()

	env, err := config.InitEnvs()
	if err != nil {
		log.Fatal("get envs failed", err)
	}

	bot := client.NewTelegramClient(env.TelegramToken)

	dbClient, err := getDbConnection(ctx, env.DbUrl, env.DbConnectionTimeout)
	if err != nil {
		log.Fatal("connection failed", err)
	}

	dbRepository := repository.NewMongoRepository(dbClient.Database(env.DbName))
	processor := event.NewProcessor(*bot, dbRepository, env.RadiusMeters, config.DefaultOffset, config.DefaultLimit)

	log.Print("service started")

	cons := consumer.NewConsumer(*processor)
	if err := cons.Start(ctx); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func getDbConnection(ctx context.Context, dbUrl string, timeOut time.Duration) (*mongo.Client, error) {
	dbClient, err := mongo.NewClient(options.Client().ApplyURI(dbUrl).SetConnectTimeout(timeOut))
	if err != nil {
		log.Fatal("db start is stopped", err)

		return nil, err
	}
	err = dbClient.Connect(ctx)
	if err != nil {
		log.Fatal("db connect is stopped", err)

		return nil, err
	}

	err = dbClient.Ping(ctx, nil)
	if err != nil {
		log.Fatal("db connect is failed: ", err)

		return nil, err
	}

	return dbClient, nil
}
