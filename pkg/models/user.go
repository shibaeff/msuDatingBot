package models

import (
	"echoBot/pkg/bot/controllers"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
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
	genderController  controllers.Controller = &controllers.GenderController{}
	photoController   controllers.Controller = &controllers.PhotoController{}
	aboutController   controllers.Controller = &controllers.AboutController{}
	facultyController controllers.Controller = controllers.NewFacultyController()
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

func (u *User) RegisterStepMessage(message *tgbotapi.Message) (reply *tgbotapi.MessageConfig, err error) {
	text := message.Text
	reply = &tgbotapi.MessageConfig{}
	reply.ParseMode = "html"
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
		errorMsg, e := facultyController.Verify(strings.ToLower(text))
		if e != nil {
			reply.Text = errorMsg
			return
		}
		u.RegiStep = regAbout
		u.Faculty = text
		reply.Text = askAbout
		return
	case regAbout:
		errorMsg, e := aboutController.Verify(text)
		if e != nil {
			reply.Text = errorMsg
			return
		}
		u.About = text
		u.RegiStep = regPhoto
		reply.Text = "Загрузите свое фото"
		return
	case regPhoto:
		errorMsg, e := photoController.Verify(message.Photo == nil &&
			message.Document == nil)
		if e != nil {
			reply.Text = errorMsg
			return
		}
		if message.Photo != nil {
			photos := *message.Photo
			photo := photos[0]
			u.PhotoLink = photo.FileID
			u.RegiStep = regOver
			return nil, nil
		} else {

		}
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

func (u *User) ReplyWithPhoto() (ret *tgbotapi.PhotoConfig) {
	ret = &tgbotapi.PhotoConfig{
		BaseFile: tgbotapi.BaseFile{
			BaseChat: tgbotapi.BaseChat{
				ChatID: u.Id,
			},
			UseExisting: true,
			FileID:      u.PhotoLink,
		},
		Caption:   u.String(),
		ParseMode: "html",
	}
	return
}

func (u *User) ChangeAbout(text string) (ret *tgbotapi.MessageConfig) {
	errorMsg, err := aboutController.Verify(text)
	if err != nil {
		return u.ReplyWithText(errorMsg)
	}
	u.About = text
	return u.ReplyWithText("Успешное изменение")
}

func (u *User) String() string {
	return fmt.Sprintf(stringify, u.Name, u.Faculty, u.Gender, u.WantGender, u.About)
}
