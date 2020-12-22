package bot

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"echoBot/pkg/models"
	"echoBot/pkg/store"
	"echoBot/pkg/timelogger"
)

const (
	waiting  = -1
	regBegin = 0
	regOver  = 7

	defaultBunchSize = 5
	noPhoto          = "none"

	timeLoggingFileName = "time.csv"

	registerCommand   = "/register"
	unseenCommand     = "/unseen"
	nextCommand       = "/next"
	usersCommand      = "/users"
	helpCommand       = "/help"
	likeCommand       = "/like"
	matchesCommand    = "/matches"
	resetCommand      = "/reset"
	profileCommand    = "/profile"
	photoCommand      = "/photo"
	startCommand      = "/start"
	facultyCommand    = "/faculty"
	aboutCommand      = "/about"
	logCommand        = "/log"
	dumpCommand       = "/dump"
	notifyCommand     = "/notify"
	reregisterCommand = "/reregister"
	feedbackCommand   = "/feedback"
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
	user, err := b.store.GetUser(message.Chat.ID)
	text := message.Text
	r := b.handleSimpleCommands(user, text)
	if r != nil {
		r.ChatID = message.Chat.ID
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
	// if registration is over
	if !user.IsReg() || user.RegiStep < regOver {
		reply, err = user.RegisterStepMessage(message)
		if err == nil {
			b.store.DeleteUser(user.Id)
			b.store.PutUser(*user)
		}
		if reply.(*tgbotapi.MessageConfig) == nil {
			// populate db for user after registration is over
			go func() {
				b.populateNotify(user)
			}()
			return user.ReplyWithPhoto(), nil
		}
		return reply, nil
	}
	if text[0] == '/' {
		split := strings.Split(text, " ")
		// in case of paired commands
	Reregister:
		if len(split) == 1 {
			switch text {
			case registerCommand:
				if !user.IsReg() {
					user.RegiStep = regBegin
					reply, _ = user.RegisterStepMessage(message)
					return
				} else {
					reply = user.ReplyWithText(alreadyRegistered)
					return
				}
			case reregisterCommand:
				b.store.DeleteUser(user.Id)
				user.RegiStep = waiting
				text = registerCommand
				goto Reregister
			case dumpCommand:
				if !b.ensureAdmin(user.UserName) {
					return user.ReplyWithText(notAdmin), nil
				}
				if !b.ensureAdmin(user.UserName) {
					return user.ReplyWithText(notAdmin), nil
				}
				b.dumpEntire()
				fileUpload := tgbotapi.NewDocumentUpload(user.Id, "dump.json")
				fileUpload.ChatID = user.Id
				fileUpload.Caption = "Ваш дамп!"
				return &fileUpload, nil
			case profileCommand:
				reply = user.ReplyWithPhoto()
				return
			case nextCommand:
				candidates, _ := b.store.GetActions().GetEvents(store.Options{
					bson.E{
						"event", store.EventUseen,
					},
				})
				candidate_id := candidates[0].Whome
				candidate, _ := b.store.GetUser(candidate_id)
				return b.replyWithCard(candidate, user.Id), nil
			case unseenCommand:
				unseen, _ := b.store.GetActions().GetEvents(store.Options{
					bson.E{
						"who", user.Id,
					},
					bson.E{
						"event", store.EventUseen,
					},
				})
				return user.ReplyWithText(string(len(unseen))), nil
			}
		}
		if len(split) == 2 {
			switch split[0] {
			case usersCommand:
				if !b.ensureAdmin(user.UserName) {
					return user.ReplyWithText(notAdmin), nil
				}
				n, err := strconv.Atoi(split[1])
				if err != nil {
					return user.ReplyWithText("Ошибка парсинга"), nil
				}
				r := b.users(n)
				r.ChatID = user.Id
				return r, nil
			case aboutCommand:
				reply = user.ChangeAbout(split[1])
				b.store.UpdUserField(user.Id, "about", user.About)
				return
			case facultyCommand:
				reply = user.ChangeFaculty(split[1])
				b.store.UpdUserField(user.Id, "faculty", user.Faculty)
				return
			case feedbackCommand:
				b.feedback(split[1])
				reply = user.ReplyWithText("Отзыв успешно доставлен")
				return
			case notifyCommand:
				if !b.ensureAdmin(user.UserName) {
					return user.ReplyWithText(notAdmin), nil
				}
				b.notifyUsers(split[1])
			case logCommand:
				if !b.ensureAdmin(user.UserName) {
					reply = user.ReplyWithText(notAdmin)
					return
				}
				if len(split) < 2 {
					reply = user.ReplyWithText("Неправильный оффсет")
					return reply, nil
				}
				var offset, err = strconv.Atoi(split[1])
				if err != nil {
					err = nil
					reply = user.ReplyWithText("Неправильный оффсет")
					return reply, nil
				}
				logs, err := b.grabLogs(offset)
				if err != nil {
					err = nil
					reply = user.ReplyWithText("Неправильный оффсет")
					return reply, nil
				}
				reply = user.ReplyWithText(logs)
				return reply, nil
			}
		}
		reply = user.ReplyWithText("Неизвестная команда")
		return
	}
	return user.ReplyWithText(notUnderstood), nil
}

func (b *bot) handleSimpleCommands(user *models.User, text string) (reply *tgbotapi.MessageConfig) {
	switch text {
	case helpCommand:
		if user != nil {
			for _, item := range b.adminsList {
				if item == user.UserName {
					return user.ReplyWithText(helpMsg + adminHelp)
				}
			}
		}
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
