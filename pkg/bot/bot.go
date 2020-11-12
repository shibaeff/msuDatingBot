package bot

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"echoBot/pkg/models"
	"echoBot/pkg/store"
)

const (
	waiting          = -1
	defaultBunchSize = 5
	registerCommand  = "/register"
	nextCommand      = "/next"
	usersCommand     = "/users"
	helpCommand      = "/help"

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
	Reply(message *tgbotapi.Message) (interface{}, error)
}

type bot struct {
	store store.Store
	api   *tgbotapi.BotAPI
}

// var Users = make(map[int64]bool)

func (b *bot) Reply(message *tgbotapi.Message) (reply interface{}, err error) {
	user, err := b.store.GetUser(message.Chat.ID)
	if err != nil {
		reply = replyWithText(greetMsg)
		err = b.store.PutUser(&models.User{
			Name:       message.Chat.FirstName,
			Faculty:    "",
			Gender:     "",
			WantGender: "",
			About:      "",
			Id:         message.Chat.ID,
			PhotoLink:  "",
			RegiStep:   waiting,
			UserName:   message.Chat.UserName,
		})
		return
	}
	if user.RegiStep != waiting && user.RegiStep < regOver {
		reply = b.registerFlow(user, message)
		return
	}
	if message.IsCommand() {
		switch message.Text {
		case helpCommand:
			reply = replyWithText(helpMsg)
			return
		case registerCommand:
			if user.RegiStep >= regOver {
				reply = replyWithText(alreadyRegistered)
				return
			}
			reply = b.registerFlow(user, message)
			// RegisterStatus[message.Chat.ID] = 1
			return
		case nextCommand:
			if user.RegiStep < regOver {
				reply = replyWithText(notRegistered)
				return
			}
			user, err = b.store.GetAny(user.Id)
			if err != nil {
				return
			}
			reply = replyWithPhoto(user)
			return
		case usersCommand:
			if user.RegiStep < regOver {
				reply = replyWithText(notRegistered)
				return
			}
			usersString, err := b.listUsers()
			if err != nil {
				return nil, err
			}
			reply = replyWithText(usersString)
			return reply, nil
		}
	}
	reply = replyWithText(notUnderstood)
	return
}

func (b *bot) listUsers() (str string, err error) {
	users, err := b.store.GetBunch(defaultBunchSize)
	if err != nil {
		log.Fatal(err)
		return
	}
	var raw []string
	for _, user := range users {
		raw = append(raw, user.String())
	}
	return strings.Join(raw, "\n"), nil
}

func replyWithText(text string) (ret *tgbotapi.MessageConfig) {
	ret = &tgbotapi.MessageConfig{
		Text: text,
	}
	ret.ReplyMarkup = menuKeyboard
	return
}

func replyWithPhoto(u *models.User) (ret *tgbotapi.PhotoConfig) {
	//ret = &tgbotapi.PhotoConfig{
	//
	//	},
	//	Caption: u.String(),
	//}
	return
}

func NewBot(store store.Store, api *tgbotapi.BotAPI) (b Bot) {
	b = &bot{store: store, api: api}
	return b
}
