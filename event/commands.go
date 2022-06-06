package event

import (
	"ArchitectureBot/models"
	"ArchitectureBot/repository"
	"context"
	"fmt"
	"strings"
)

const (
	StartCmd = "/start"
	HelpCmd  = "/help"
	//Todo убрать
	SetGeoCmd = "/setpos"
)

func (p *Processor) doCommand(ctx context.Context, text string, location *models.Location, chatId int64) error {
	text = strings.TrimSpace(text)

	if location != nil {
		return p.SendBuilds(ctx, chatId, location)
	}

	switch text {
	case StartCmd:
		return p.SendStart(ctx, chatId)
	case HelpCmd:
		return p.SendHelp(ctx, chatId)
	case SetGeoCmd:
		return p.SetBuilds(ctx)
	default:
		return p.SendUnknown(ctx, chatId)
	}
}

func (p *Processor) SendHelp(ctx context.Context, chatId int64) error {
	return p.client.SendMessage(chatId, "Инструкция")
}

func (p *Processor) SendStart(ctx context.Context, chatId int64) error {
	return p.client.SendMessage(chatId, "Приветсвую тебя, мой господин")
}

func (p *Processor) SendUnknown(ctx context.Context, chatId int64) error {
	return p.client.SendMessage(chatId, "Хм, эта команда мне не знакома")
}

func (p *Processor) SendBuilds(ctx context.Context, chatId int64, location *models.Location) error {
	builds, err := p.getBuilds(ctx, location.Longitude, location.Latitude, p.radius)
	if err != nil {
		return err
	}

	var text string
	if len(builds) == 0 {
		text = "Ничего не найдено поблизости"

		return p.client.SendMessage(chatId, text)
	}

	for _, v := range builds {
		text += p.createBuildMsg(v)
	}

	return p.client.SendMessage(chatId, text)
}

func (p *Processor) getBuilds(ctx context.Context, longitude, latitude string, radius int) ([]repository.ArchitectBuilds, error) {
	builds, err := p.db.GetBuildsByDistance(ctx, longitude, latitude, radius)
	if err != nil {
		return nil, err
	}

	return builds, nil
}

func (p *Processor) createBuildMsg(build repository.ArchitectBuilds) string {
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
