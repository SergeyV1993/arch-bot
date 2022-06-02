package event

import (
	"ArchitectureBot/client"
	"ArchitectureBot/models"
	"ArchitectureBot/repository"
	"context"
	"github.com/pkg/errors"
)

type TelegramClientInterface interface {
	Updates(offset, limit int) ([]models.Update, error)
	SendMessage(chatId int64, text string) error
}

type MongoRepositoryInterface interface {
	SetBuilds(ctx context.Context, build repository.ArchitectBuilds) error
	GetBuildsByDistance(ctx context.Context, longitude, latitude string, distance int) ([]repository.ArchitectBuilds, error)
}

type Processor struct {
	client TelegramClientInterface
	db     MongoRepositoryInterface

	offset int
	radius int
}

func NewProcessor(client client.TelegramClient, db MongoRepositoryInterface, offset, radius int) *Processor {
	return &Processor{
		client: &client,
		db:     db,
		offset: offset,
		radius: radius,
	}
}

func (p *Processor) Fetch(ctx context.Context, limit int) ([]models.Update, error) {
	updates, err := p.client.Updates(p.offset, limit)
	if err != nil {
		return nil, errors.Wrap(err, "can't get events")
	}

	if len(updates) == 0 {
		return nil, nil
	}

	p.offset = updates[len(updates)-1].ID + 1

	return updates, nil
}

func (p *Processor) Process(ctx context.Context, event models.Update) error {
	if err := p.doCommand(ctx, event.Message.Text, event.Message.Location, event.Message.ChatID); err != nil {
		return errors.Wrap(err, "can't process message")
	}

	return nil
}