package consumer

import (
	"ArchitectureBot/event"
	"ArchitectureBot/models"
	"context"
	"log"
	"time"
)

type ProcessorInterface interface {
	Fetch(ctx context.Context) ([]models.Update, error)
	Process(ctx context.Context, event models.Update) error
}

type Consumer struct {
	processorService ProcessorInterface
}

func NewConsumer(processorService event.Processor) *Consumer {
	return &Consumer{
		processorService: &processorService,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		gotEvents, err := c.processorService.Fetch(ctx)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(ctx, gotEvents); err != nil {
			log.Print(err)

			continue
		}
	}
}

func (c *Consumer) handleEvents(ctx context.Context, updates []models.Update) error {
	for _, update := range updates {
		log.Printf("got new event: %s", update.Message.Text)

		if err := c.processorService.Process(ctx, update); err != nil {
			log.Printf("can't handle event: %s", err.Error())

			continue
		}
	}

	return nil
}
