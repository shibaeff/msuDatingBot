package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"echoBot/pkg/models"
	"echoBot/pkg/store"
)

const (
	waiting         = -1
	registerCommand = "/register"
	nextCommand     = "/next"
	usersCommand    = "/users"
	helpCommand     = "/help"

	greetMsg          = "Добро пожаловать в бота знакомств. Начните с /register."
	notUnderstood     = "Пожалуйста, выберите действие из меню"
	alreadyRegistered = "Вы уже зарегистрированы!"
	notRegistered     = "Вы не зарегистрированы!"

	helpMsg = "Бот знакомств поможет Вам найти интересных людей. \n " +
		"/register - регистрация\n" +
		"/next - получить следующее предложение\n" +
		"/users - вывести список активных пользователей \n"
)

var (
	registerButton = tgbotapi.KeyboardButton{Text: registerCommand}
	helpButton     = tgbotapi.KeyboardButton{Text: helpCommand}
	nextButton     = tgbotapi.KeyboardButton{Text: nextCommand}
	usersButton    = tgbotapi.KeyboardButton{Text: usersCommand}
	menuButtons    = []tgbotapi.KeyboardButton{registerButton, helpButton, nextButton, usersButton}
	menuKeyboard   = tgbotapi.NewReplyKeyboard(menuButtons)
)

type Bot interface {
	Reply(message *tgbotapi.Message) *tgbotapi.MessageConfig
}

type bot struct {
	store store.Store
}

func replyWithText(text string) (ret *tgbotapi.MessageConfig) {
	ret = &tgbotapi.MessageConfig{
		Text: text,
	}
	ret.ReplyMarkup = menuKeyboard
	return
}

// var Users = make(map[int64]bool)
var RegisterStatus = make(map[int64]int64)
var Photos = make(map[int64]string)

func (b *bot) Reply(message *tgbotapi.Message) (reply *tgbotapi.MessageConfig) {
	_, err := b.store.GetUser(message.Chat.ID)
	if err != nil {
		reply = replyWithText(greetMsg)
		b.store.PutUser(&models.User{
			Name:       "",
			Faculty:    "",
			Gender:     "",
			WantGender: "",
			About:      "",
			Id:         message.Chat.ID,
			PhotoLink:  "",
			RegiStep:   0,
			UserName:   message.Chat.UserName,
		})
		return
	}
	_, ok = RegisterStatus[message.Chat.ID]
	if ok && RegisterStatus[message.Chat.ID] < regOver {
		reply = registerFlow(message)
		return
	}
	if message.IsCommand() {
		switch message.Text {
		case helpCommand:
			reply = replyWithText(helpMsg)
			return
		case registerCommand:
			if RegisterStatus[message.Chat.ID] >= regOver {
				reply = replyWithText(alreadyRegistered)
				return
			}
			reply = registerFlow(message)
			// RegisterStatus[message.Chat.ID] = 1
			return
		case nextCommand:
			if RegisterStatus[message.Chat.ID] < regOver {
				reply = replyWithText(notRegistered)
				return
			}
		case usersCommand:
			if RegisterStatus[message.Chat.ID] < regOver {
				reply = replyWithText(notRegistered)
				return
			}
		}
	}
	reply = replyWithText(notUnderstood)
	return
}

func NewBot(store store.Store) (b Bot) {
	b = &bot{store: store}
	return b
}
