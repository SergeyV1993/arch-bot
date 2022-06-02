package event

import (
	"ArchitectureBot/models"
	"ArchitectureBot/repository"
	"context"
	"fmt"
	"strings"
)

const (
	StartCmd  = "/start"
	HelpCmd   = "/help"
	SetGeoCmd = "/setpos"
)

func (p *Processor) doCommand(ctx context.Context, text string, location *models.Location, chatId int64) error {
	text = strings.TrimSpace(text)

	if location != nil {
		return p.GetBuilds(ctx, chatId, location, p.radius)
	}

	switch text {
	case StartCmd:
		return p.SendStart(chatId)
	case HelpCmd:
		return p.SendHelp(chatId)
	case SetGeoCmd:
		return p.SetBuilds(ctx)
	default:
		return p.SendUnknown(chatId)
	}
}

func (p *Processor) SendHelp(chatId int64) error {
	return p.client.SendMessage(chatId, "Инструкция")
}

func (p *Processor) SendStart(chatId int64) error {
	return p.client.SendMessage(chatId, "Приветсвую тебя, мой господин")
}

func (p *Processor) SendUnknown(chatId int64) error {
	return p.client.SendMessage(chatId, "Не знаю, что сказать")
}

func (p *Processor) GetBuilds(ctx context.Context, chatId int64, location *models.Location, radius int) error {
	builds, err := p.db.GetBuildsByDistance(ctx, location.Longitude, location.Latitude, radius)
	if err != nil {
		return err
	}

	var text string
	if len(builds) == 0 {
		text = "Ничего не найдено поблизости"

		return p.client.SendMessage(chatId, text)
	}

	for _, v := range builds {
		text += p.createMsg(v)
	}

	return p.client.SendMessage(chatId, text)
}

func (p *Processor) createMsg(build repository.ArchitectBuilds) string {
	text := fmt.Sprintf("*%s*", build.Name) + "\n" +
		fmt.Sprintf("Приблизительное расстояние к достопримечательности: _%.1f_ м.", build.Distance) + "\n" +
		fmt.Sprintf("[%s](%s)", build.Address, build.LinkMapAddress) + "\n" +
		fmt.Sprintf("*Описание: *_%s_", build.Description) + "\n" +
		fmt.Sprintf("[Более подробно по ссылке](%s)", build.Link) +
		"\n" + "\n"

	return text
}

//TODO debug
func (p *Processor) SetBuilds(ctx context.Context) error {
	err := p.db.SetBuilds(ctx, repository.ArchitectBuilds{})
	if err != nil {
		return err
	}

	return nil
}
