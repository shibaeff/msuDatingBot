package bot

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"echoBot/pkg/models"
)

const (
	regBegin      = 0
	regName       = 1
	regGender     = 2
	regWantGender = 3
	regFaculty    = 4
	regAbout      = 5
	regPhoto      = 6
	regOver       = 7

	askName       = "Пожалуйста, введите свое имя"
	askGender     = "Пожалуйста, введите Ваш пол м/ж"
	askWantGender = "Кого ищем: м/ж/любой?"
	askFaculty    = "С какого Вы факультета?"
	askAbout      = "Напишите немного о себе"
)

func registerStep(text string) (reply *tgbotapi.MessageConfig) {
	reply = &tgbotapi.MessageConfig{
		Text: text,
	}
	return
}

func (b *bot) updateRegStatus(id, status int64) error {
	return b.store.UpdUserField(id, "registep", status)
}

func (b *bot) handleGender(user *models.User, query *tgbotapi.CallbackQuery) {
	b.api.Send(tgbotapi.EditMessageTextConfig{

	})
}

func (b *bot) registerFlow(user *models.User, text string) (reply *tgbotapi.MessageConfig) {
	switch user.RegiStep {
	case waiting:


	return
}
