package bot

import (
	"context"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"echoBot/pkg/models"
	"echoBot/pkg/store"
	"echoBot/pkg/timelogger"
)

const (
	waiting  = -1
	regBegin = 0

	defaultBunchSize = 5
	noPhoto          = "none"

	timeLoggingFileName = "time.csv"

	registerCommand   = "/register"
	nextCommand       = "/next"
	usersCommand      = "/users"
	helpCommand       = "/help"
	likeCommand       = "/like"
	matchesCommand    = "/matches"
	resetCommand      = "/reset"
	profileCommand    = "/profile"
	photoCommand      = "/photo"
	startCommand      = "/start"
	cancelCommand     = "/cancel"
	facultyCommand    = "/faculty"
	aboutCommand      = "/about"
	logCommand        = "/log"
	dumpCommand       = "/dump"
	notifyAll         = "/notify"
	reregisterCommand = "/reregister"
	feedbackCommand   = "/feedback"
	numbers           = "/numbers"
	deleteCommand     = "/delete"
	pauseCommand      = "/pause"

	greetMsg          = "Привет! ✨\nЭто бот знакомств МГУ. Работает аналогично Тиндеру 😉\n\nДля регистрации вызывай: /register, для отмены: /cancel. Бот запросит имя, фоточку и пару слов о себе.\n\nПредложения и баги пишите в /feedback."
	notUnderstood     = "Пожалуйста, выберите действие из меню"
	alreadyRegistered = "Вы уже зарегистрированы!"
	notRegistered     = "Вы не зарегистрированы!"
	notAdmin          = "Вы не админ"
	pleaseSendAgain   = "Пожалуйста, сделайте запрос еще раз"
)

type Bot interface {
	ReplyMessage(ctx context.Context, message *tgbotapi.Message) (interface{}, error)
	HandleCallbackQuery(ctx context.Context, query *tgbotapi.CallbackQuery) (interface{}, error)
	GetStore() store.Store
}

type bot struct {
	store       store.Store
	api         *tgbotapi.BotAPI
	logFile     *os.File
	timeloggers map[string]timelogger.TimeLogger
	adminsList  []string
	actionsLog  *log.Logger
}

func (b *bot) GetStore() store.Store {
	return b.store
}

func (b *bot) HandleCallbackQuery(context context.Context, query *tgbotapi.CallbackQuery) (reply interface{}, err error) {
	user, err := b.store.GetUser(int64(query.From.ID))
	b.api.AnswerCallbackQuery(tgbotapi.CallbackConfig{
		CallbackQueryID: query.ID,
		Text:            "",
		ShowAlert:       false,
		URL:             "",
		CacheTime:       0,
	})
	if !user.IsReg() {
		reply = user.RegisterStepInline(query)
		b.store.DeleteUser(user.Id)
		b.store.PutUser(*user)
		return
	}
	return nil, nil
}

func (b *bot) ReplyMessage(context context.Context, message *tgbotapi.Message) (reply interface{}, err error) {
	// fast track
	switch message.Text {
	case deleteCommand:
		return b.deleteUser(message.Chat.ID), nil
	}
	user := &models.User{}
	text := message.Text
	r := handleSimpleCommands(user, text)
	if r != nil {
		reply.(*tgbotapi.MessageConfig).ChatID = message.Chat.ID
		return r, nil
	}
	// check if user is registered
	// unregistered users are allowed only to do /start, /help, /register
	user, err = b.store.GetUser(message.Chat.ID)
	// Putting user in the db
	if err != nil {
		u := models.User{
			Name:       message.Chat.FirstName,
			Faculty:    "",
			Gender:     "",
			WantGender: "",
			About:      "",
			Id:         message.Chat.ID,
			PhotoLink:  "",
			RegiStep:   waiting,
			UserName:   message.Chat.UserName,
		}
		b.store.PutUser(u)
		user = &u
	}
	if !user.IsReg() {
		reply, err = user.RegisterStepMessage(text)
		if err == nil {
			b.store.DeleteUser(user.Id)
			b.store.PutUser(*user)
		}
		return reply, nil
	}
	if text[0] == '/' {
		split := strings.Split(text, " ")
		// in case of paired commands
		if len(split) == 1 {
			switch text {
			case registerCommand:
				user.RegiStep = regBegin
				// b.store.UpdUserField(user.Id, "registep", user.RegiStep)
				reply, _ = user.RegisterStepMessage(text)
				return
			}
		}
		if len(split) == 2 {

		}
		reply = user.ReplyWithText("Неизвестная команда")
		return
	}
	return user.ReplyWithText(notUnderstood), nil
}

func handleSimpleCommands(user *models.User, text string) (reply *tgbotapi.MessageConfig) {
	switch text {
	case helpCommand:
		return user.ReplyWithText(helpMsg)
	case startCommand:
		return user.ReplyWithText(greetMsg)
	}
	return nil
}
func (b *bot) setTimeLoggers() {
	b.timeloggers = make(map[string]timelogger.TimeLogger)
	b.timeloggers[startCommand] = timelogger.NewTimeLogger(startCommand, timeLoggingFileName)
}

func (b *bot) setActionLoggers() {
	file, err := os.OpenFile("actions.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("Cannot create or open log file")
	}
	b.actionsLog = log.New(file, "Common Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
}

func NewBot(store store.Store, api *tgbotapi.BotAPI, logFile *os.File, admins []string) (b Bot) {
	b = &bot{
		store:      store,
		api:        api,
		logFile:    logFile,
		adminsList: admins,
	}
	b.(*bot).setTimeLoggers()
	b.(*bot).setActionLoggers()
	return b
}
