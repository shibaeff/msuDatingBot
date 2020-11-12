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
	likeCommand      = "/like"
	matchesCommand   = "/matches"

	greetMsg          = "Добро пожаловать в бота знакомств. Начните с /register."
	notUnderstood     = "Пожалуйста, выберите действие из меню"
	alreadyRegistered = "Вы уже зарегистрированы!"
	notRegistered     = "Вы не зарегистрированы!"

	helpMsg = "🔍 Поиск:\n• /next — просмотреть следующую анкету\n• /matches — взаимные лайки\n• /info — посмотреть информацию\n\n📋 Профиль:\n• /profile — посмотреть как выглядит ваш профиль\n• /register — пройти регистрацию заново \n• /photo — обновить фото \n• /faculty — обновить факультет\n• /about — обновить описание \n• /settings — прочие настройки профиля\n\n⚙️ Прочие команды:\n• /start — общее описание бота\n• /help — вызов этого сообщения\n• /cancel — отмена текущей операции\n• /reset — сбросить все свои оценки (аккуратно!)"
)

var (
	registerButton = tgbotapi.KeyboardButton{Text: registerCommand}
	helpButton     = tgbotapi.KeyboardButton{Text: helpCommand}
	nextButton     = tgbotapi.KeyboardButton{Text: nextCommand}
	usersButton    = tgbotapi.KeyboardButton{Text: likeCommand}
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
			return
		case nextCommand:
			if user.RegiStep < regOver {
				reply = replyWithText(notRegistered)
				return
			}
			newuser, e := b.store.GetAny(user.Id)
			if e != nil {
				reply = replyWithText("Не можем подобрать вариант")
				return
			}
			e = b.store.PutSeen(user.Id, newuser.Id)
			if e != nil {
				reply = replyWithText("Не можем подобрать вариант")
				return
			}
			reply = replyWithPhoto(newuser, message.Chat.ID)
			return
		case likeCommand:
			entry, e := b.store.GetSeen(user.Id)
			if e != nil {
				reply = replyWithText("failed to put your like")
				return
			}
			likee := entry.Whome[len(entry.Whome)-1]
			e = b.store.PutLike(user.Id, likee)
			if e != nil {
				reply = replyWithText("failed to put your like")
				return
			}
			reply = replyWithText("Успешный лайк!")
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

func NewBot(store store.Store, api *tgbotapi.BotAPI) (b Bot) {
	b = &bot{store: store, api: api}
	return b
}
