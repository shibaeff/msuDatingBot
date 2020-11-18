package bot

import (
	"errors"
	"fmt"
	"strings"

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

func (b *bot) prepareMatches(userId int64) (resp string, err error) {
	entry, err := b.store.GetMatchesRegistry().GetList(userId)
	if err != nil {
		return "Матчей нет", nil
	}
	matches := entry.Whome
	if len(matches) == 0 {
		return "Матчей нет", nil
	}
	raw := []string{}
	for _, match := range matches {
		user, err := b.store.GetUser(match)
		if err != nil {
			continue
		}
		raw = append(raw, fmt.Sprintf("@%s\n", user.UserName))
	}
	resp = matchesList + strings.Join(raw, "")
	return
}

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

func find(slice []int64, val int64) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
