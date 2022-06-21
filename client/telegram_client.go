package client

import (
	"ArchitectureBot/models"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/pkg/errors"
	"log"
)

type TelegramClient struct {
	Client *tgbotapi.BotAPI
}

func NewTelegramClient(telegramToken string) *TelegramClient {
	client, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Fatalf("I don't parse telegram token, err= %s", err)
	}

	return &TelegramClient{Client: client}
}

func (tc *TelegramClient) Updates(offset, limit int) ([]models.Update, error) {
	uc := tgbotapi.UpdateConfig{
		Offset:  offset,
		Limit:   limit,
		Timeout: 60,
	}

	updates, err := tc.Client.GetUpdates(uc)
	if err != nil {
		return nil, errors.Wrap(err, "Can't get updates channel")
	}

	res := make([]models.Update, 0, len(updates))

	for _, update := range updates {
		convertedUpdate, err := tc.convertUpdate(update)
		if err != nil {
			log.Print("can't convert update to message")
			continue
		}

		res = append(res, *convertedUpdate)
	}

	return res, nil
}

func (tc *TelegramClient) SendMessage(chatId int64, text string) error {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = "markdown"
	msg.DisableWebPagePreview = true

	locationBtn := tgbotapi.NewKeyboardButtonLocation("Отправить локацию")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{locationBtn})

	_, err := tc.Client.Send(msg)
	if err != nil {
		return errors.Wrap(err, "Can't send message")
	}

	return nil
}

func (tc *TelegramClient) convertUpdate(upd tgbotapi.Update) (*models.Update, error) {
	convertedUpd := models.Update{}
	convertedUpd.ID = upd.UpdateID
	convertedUpd.Message = &models.IncomingMessage{}

	if upd.Message != nil {
		convertedUpd.Message.Text = upd.Message.Text
		convertedUpd.Message.Username = upd.Message.Chat.UserName
		convertedUpd.Message.ChatID = upd.Message.Chat.ID
		if upd.Message.Location != nil {
			convertedUpd.Message.Location = &models.Location{
				Longitude: fmt.Sprintf("%f", upd.Message.Location.Longitude),
				Latitude:  fmt.Sprintf("%f", upd.Message.Location.Latitude),
			}
		}
	}

	return &convertedUpd, nil
}
