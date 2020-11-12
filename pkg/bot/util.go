package bot

import (
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"echoBot/pkg/models"
)

func (b *bot) parseLikee(message *tgbotapi.Message) (id int64, err error) {
	if message.ReplyToMessage == nil {
		return -1, errors.New("nothing to reply to")
	}
	text := message.ReplyToMessage.Text
	_, err = fmt.Scanf(text, &id)
	if err != nil {
		return -1, err
	}
	return
}

// TODO
//func (b *bot) prepareMatches(userId int64) (err error) {
//
//}

func replyWithText(text string) (ret *tgbotapi.MessageConfig) {
	ret = &tgbotapi.MessageConfig{
		Text: text,
	}
	ret.ReplyMarkup = menuKeyboard
	return
}

func replyWithPhoto(u *models.User, to int64) (ret *tgbotapi.PhotoConfig) {
	ret = &tgbotapi.PhotoConfig{
		BaseFile: tgbotapi.BaseFile{
			BaseChat: tgbotapi.BaseChat{
				ChatID: to,
			},
			UseExisting: true,
			FileID:      u.PhotoLink,
		},
		Caption: u.String(),
	}
	return
}
