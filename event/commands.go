package event

import (
	"ArchitectureBot/event/command_msg"
	"ArchitectureBot/models"
	"context"
	"strings"
)

const (
	HelpCmd = "/help"
)

func (p *Processor) doCommand(ctx context.Context, upd models.Update) error {
	if upd.Message == nil {
		return nil
	}

	text := strings.TrimSpace(upd.Message.Text)

	if upd.Message.Location != nil {
		return p.SendBuilds(ctx, upd.Message.ChatID, upd.Message.Location)
	}

	switch text {
	case HelpCmd:
		return p.SendHelp(ctx, upd.Message.ChatID, upd.Message.Username)
	default:
		return p.SendUnknown(ctx, upd.Message.ChatID)
	}
}

func (p *Processor) SendHelp(ctx context.Context, chatId int64, username string) error {
	return p.client.SendMessage(chatId, command_msg.CreateSetupMsg(username))
}

func (p *Processor) SendUnknown(ctx context.Context, chatId int64) error {
	return p.client.SendMessage(chatId, command_msg.CreateUnknownMsg())
}

func (p *Processor) SendBuilds(ctx context.Context, chatId int64, location *models.Location) error {
	builds, err := p.db.GetBuildsByDistance(ctx, location.Longitude, location.Latitude, p.radius)
	if err != nil {
		return err
	}

	var text string
	if len(builds) == 0 {
		text = "Ничего не найдено поблизости"

		return p.client.SendMessage(chatId, text)
	}

	for _, v := range builds {
		text += command_msg.CreateBuildMsg(v)
	}

	return p.client.SendMessage(chatId, text)
}
