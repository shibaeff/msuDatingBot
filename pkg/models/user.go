package models

import (
	"echoBot/pkg/bot/controllers"
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
	askGender     = "Пожалуйста, введите Ваш пол"
	askWantGender = "Кого ищем: м/ж/любой?"
	askFaculty    = "С какого Вы факультета?"
	askAbout      = "Напишите немного о себе"
)

// controllers
var (
	genderController controllers.Controller = &controllers.GenderController{}
	photoController  controllers.Controller = &controllers.PhotoController{}
	aboutController  controllers.Controller = &controllers.AboutController{}
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

func (u *User) RegisterStepMessage(text string) (reply *tgbotapi.MessageConfig, err error) {
	reply = &tgbotapi.MessageConfig{}
	reply.ChatID = u.Id
	switch u.RegiStep {
	case regBegin, regWaiting:
		u.RegiStep = regName
		reply.Text = askName
		return
	case regName:
		u.RegiStep = regGender
		u.Name = text
		reply.Text = askGender
		reply.ReplyMarkup = genderKeyboard
		return
	case regFaculty:
		u.RegiStep = regAbout
		u.Faculty = text // TODO add controller
		reply.Text = askAbout
		return
	case regAbout:
		errorMsg, err := aboutController.Verify(text)
		if err != nil {
			reply.Text = errorMsg
			return
		}
		u.RegiStep = regOver
		reply.Text = u.String()
		return
	}
	return u.ReplyWithText("Пожалуйста, следуйте подсказкам бота!"), nil
}

func (u *User) IsReg() bool {
	return u.RegiStep >= regOver
}

func (u *User) RegisterStepInline(q *tgbotapi.CallbackQuery) (reply *tgbotapi.MessageConfig) {
	switch u.RegiStep {
	case regGender:
		warning, err := genderController.Verify(q.Data)
		if err != nil {
			return u.ReplyWithText(warning)
		}
		u.RegiStep = regWantGender
		u.Gender = q.Data
		reply = u.ReplyWithText(askWantGender)
		reply.ReplyMarkup = wantGenderKeyboard
		return
	case regWantGender:
		warning, err := genderController.Verify(q.Data)
		if err != nil {
			return u.ReplyWithText(warning)
		}
		u.RegiStep = regFaculty
		u.WantGender = q.Data
		reply = u.ReplyWithText(askFaculty)
		reply = u.ReplyWithText(askFaculty)
		return
	}
	return u.ReplyWithText("Пожалуйста, следуйте подсказкам бота!")
}

func (u *User) String() string {
	return fmt.Sprintf(stringify, u.Name, u.Faculty, u.Gender, u.WantGender, u.About)
}
