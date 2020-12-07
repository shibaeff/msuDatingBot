package models

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	stringify = "<strong>Имя:</strong> %s\n" +
		"<strong>Факультет:</strong> %s\n" +
		"<strong>Пол:</strong> %s\n" +
		"<strong>Пол собеседника:</strong> %s\n" +
		"<strong>О себе:</strong> %s\n"

	nextCommand    = "/next"
	helpCommand    = "/help"
	matchesCommand = "/matches"
	profileCommand = "/profile"

	regWaiting    = -1
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

var (
	profileButton = tgbotapi.KeyboardButton{Text: profileCommand}
	helpButton    = tgbotapi.KeyboardButton{Text: helpCommand}
	matchesButton = tgbotapi.KeyboardButton{Text: matchesCommand}
	nextButton    = tgbotapi.KeyboardButton{Text: nextCommand}
	menuButtons   = []tgbotapi.KeyboardButton{profileButton, helpButton, matchesButton, nextButton}
	menuKeyboard  = tgbotapi.NewReplyKeyboard(menuButtons)
)

type User struct {
	Name       string
	Faculty    string
	Gender     string
	WantGender string
	About      string
	Id         int64
	PhotoLink  string
	RegiStep   int64
	UserName   string
}

func (u *User) ReplyWithText(text string) (ret *tgbotapi.MessageConfig) {
	ret = &tgbotapi.MessageConfig{
		Text: text,
	}
	ret.ReplyMarkup = menuKeyboard
	ret.ChatID = u.Id
	return
}

func (u *User) RegisterStepMessage(text string) (reply *tgbotapi.MessageConfig) {
	reply = &tgbotapi.MessageConfig{}
	reply.ChatID = u.Id
	switch u.RegiStep {
	case regWaiting:
		u.RegiStep = regBegin
		reply.Text = askName
		return
	}
	return
}

func (u *User) IsReg() bool {
	return u.RegiStep >= regOver
}

func RegisterStepInline(q *tgbotapi.CallbackQuery) (reply *tgbotapi.MessageConfig) {
	return
}
func (u *User) String() string {
	return fmt.Sprintf(stringify, u.Name, u.Faculty, u.Gender, u.WantGender, u.About)
}
