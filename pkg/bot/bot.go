package bot

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"echoBot/pkg/bot/controllers"
	"echoBot/pkg/models"
	"echoBot/pkg/store"
	"echoBot/pkg/timelogger"
)

const (
	waiting          = -1
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
	dumpCommand       = "/dump"
	notifyAll         = "/notify"
	reregisterCommand = "/reregister"
	feedbackCommand   = "/feedback"

	greetMsg          = "Добро пожаловать в бота знакомств. Начните с /register."
	notUnderstood     = "Пожалуйста, выберите действие из меню"
	alreadyRegistered = "Вы уже зарегистрированы!"
	notRegistered     = "Вы не зарегистрированы!"
	notAdmin          = "Вы не админ"
	pleaseSendAgain   = "Пожалуйста, сделайте запрос еще раз"
)

var (
	registerButton = tgbotapi.KeyboardButton{Text: registerCommand}
	helpButton     = tgbotapi.KeyboardButton{Text: helpCommand}
	matchesButton  = tgbotapi.KeyboardButton{Text: matchesCommand}
	nextButton     = tgbotapi.KeyboardButton{Text: nextCommand}
	menuButtons    = []tgbotapi.KeyboardButton{registerButton, helpButton, matchesButton, nextButton}
	menuKeyboard   = tgbotapi.NewReplyKeyboard(menuButtons)
)

type Bot interface {
	Reply(message *tgbotapi.Message) (interface{}, error)
	GetStore() store.Store
}

type bot struct {
	store            store.Store
	api              *tgbotapi.BotAPI
	genderController controllers.Controller
	photoController  controllers.Controller
	logFile          *os.File
	timeloggers      map[string]timelogger.TimeLogger
	adminsList       []string
}

func (b *bot) GetStore() store.Store {
	return b.store
}

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
		if message.Text == cancelCommand {
			if user.RegiStep != regPhoto {
				b.store.DeleteUser(user.Id)
				reply = replyWithText("Откат регистрации")
				return
			} else {
				b.store.UpdUserField(user.Id, "registep", regOver)
				reply = replyWithText("Откатываемся к старой информации")
				return
			}
		}
		reply = b.registerFlow(user, message)
		if user.RegiStep == regOver {
			users, _ := b.store.GetAllUsers()
			for _, cur_user := range users {
				if b.ensureGender(user, cur_user) {
					b.store.PutUnseen(user.Id, cur_user.Id)
					b.store.PutUnseen(cur_user.Id, user.Id)
				}
			}
		}
		return
	}
	if message.Text[0] == '/' || message.Text == likeEmoji || message.Text == nextEmoji {
		split := strings.Split(message.Text, " ")
		switch split[0] {
		case reregisterCommand:
			user.RegiStep = regName
			b.store.UpdUserField(user.Id, "registep", user.RegiStep)
			return registerStep("Начинаем новую регистрацию!" + askName), nil
		case notifyAll:
			if !b.ensureAdmin(user.UserName) {
				reply = replyWithText(notAdmin)
				return
			}
			reply, _ = b.notifyUsers(split[1])
			return
		case feedbackCommand:
			repl := replyWithText(message.Text)
			adminmsg := &MessageToAdmin{
				repl,
				b.adminsList,
			}
			return adminmsg, nil
		case dumpCommand:
			if !b.ensureAdmin(user.UserName) {
				reply = replyWithText(notAdmin)
				return
			}
			if len(split) < 2 {
				reply = replyWithText("Неправильный оффсет")
				return reply, nil
			}
			var offset, err = strconv.Atoi(split[1])
			if err != nil {
				err = nil
				reply = replyWithText("Неправильный оффсет")
				return reply, nil
			}
			logs, err := b.grabLogs(offset)
			if err != nil {
				err = nil
				reply = replyWithText("Неправильный оффсет")
				return reply, nil
			}
			reply = replyWithText(logs)
			return reply, nil
		case aboutCommand:
			about := strings.Split(message.Text, " ")[1]
			err = b.store.UpdUserField(user.Id, "about", about)
			if err != nil {
				reply = replyWithText("Ошибка обновления!")
				return
			}
			reply = replyWithText(fmt.Sprintf("Обновили информацию на %s", about))
			return
		case facultyCommand:
			faculty := strings.Split(message.Text, " ")[1]
			err = b.store.UpdUserField(user.Id, "faculty", faculty)
			if err != nil {
				reply = replyWithText("Ошибка обновления!")
				return
			}
			reply = replyWithText(fmt.Sprintf("Обновили факультет на %s", faculty))
			return
		case startCommand:
			b.timeloggers[startCommand].Start()
			reply = replyWithText(greetMsg)
			b.timeloggers[startCommand].End()
			return
		case helpCommand:
			if b.ensureAdmin(user.UserName) {
				reply = replyWithText(helpMsg + adminHelp)
				return
			}
			reply = replyWithText(helpMsg)
			return
		case registerCommand:
			if user.RegiStep >= regOver {
				reply = replyWithText(alreadyRegistered)
				return
			}
			reply = b.registerFlow(user, message)
			return
		case "/unseen":
			unseen, _ := b.store.GetUnseen(user.Id)
			ret := []string{}
			for _, item := range unseen {
				ret = append(ret, strconv.Itoa(int(item.Whome)))
			}
			reply = replyWithText(strings.Join(ret, " "))
			return
		case nextCommand, nextEmoji:
			if user.RegiStep < regOver {
				reply = replyWithText(notRegistered)
				return
			}
			unseen, e := b.store.GetUnseen(user.Id)
			if len(unseen) == 0 || e != nil {
				reply = replyWithText("Вы просмотрели всех пользователей на данный момент")
				return reply, nil
			}
			unseen_user, _ := b.store.GetUser(unseen[0].Whome)
			b.store.GetUnseenRegistry().DeleteItem(user.Id, unseen_user.Id)
			b.store.GetSeenRegistry().AddToList(user.Id, unseen_user.Id)
			reply = replyWithCard(unseen_user, user.Id)
			return
		case likeCommand, likeEmoji:
			entry, e := b.store.GetSeen(user.Id)
			if e != nil {
				reply = replyWithText("failed to put your like")
				return
			}
			likee := entry[len(entry)-1].Whome
			e = b.store.PutLike(user.Id, likee)
			if e != nil {
				reply = replyWithText("failed to put your like")
				return
			}
			reply = replyWithText("Успешный лайк!")
			likee_entries, err := b.store.GetLikes(likee)
			if err == nil {
				_, ok1 := find(likee_entries, user.Id)
				if ok1 {
					user_entry, err := b.store.GetLikes(user.Id)
					if err != nil {
						return reply, nil
					}
					_, ok1 = find(user_entry, likee)
					if ok1 {
						likee_user, err := b.store.GetUser(likee)
						if err != nil {
							reply = replyWithText("Такого пользователя уже нет")
						}
						reply1 := replyWithText(fmt.Sprintf(matchMsg, likee_user.UserName))
						reply1.ChatID = user.Id
						reply2 := replyWithText(fmt.Sprintf(matchMsg, user.UserName))
						reply2.ChatID = likee_user.Id
						reply = &Match{
							Msg1: reply1,
							Msg2: reply2,
						}
						if !b.store.GetMatchesRegistry().IsPresent(user.Id, likee_user.Id) {
							b.store.GetMatchesRegistry().AddToList(user.Id, likee_user.Id)
						}
						if !b.store.GetMatchesRegistry().IsPresent(likee_user.Id, user.Id) {
							b.store.GetMatchesRegistry().AddToList(likee_user.Id, user.Id)
						}
						return reply, nil
					} else {
						return reply, nil
					}
				} else {
					return reply, nil
				}

			}
			return reply, nil

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
		case matchesCommand:
			matches, _ := b.prepareMatches(user.Id)
			reply = replyWithText(matches)
			return
		case resetCommand:
			b.store.DeleteFromRegistires(user.Id)
			reply = replyWithText("Ваши оценки сброшены!")
			return
		case profileCommand:
			reply = replyWithPhoto(user, message.Chat.ID)
			return
		case photoCommand:
			err = b.store.UpdUserField(user.Id, "photolink", noPhoto)
			if err != nil {
				reply = replyWithText("Ошибка обновления фото")
				return
			}
			err = b.store.UpdUserField(user.Id, "registep", regPhoto)
			if err != nil {
				reply = replyWithText("Ошибка обновления фото")
				return
			}
			reply = replyWithText("Ждем ваше фото!")
			return
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
		raw = append(raw, fmt.Sprintf("@%s", user.UserName))
	}
	return strings.Join(raw, "\n"), nil
}

func (b *bot) setTimeLoggers() {
	b.timeloggers = make(map[string]timelogger.TimeLogger)
	b.timeloggers[startCommand] = timelogger.NewTimeLogger(startCommand, timeLoggingFileName)
}

func NewBot(store store.Store, api *tgbotapi.BotAPI, logFile *os.File, admins []string) (b Bot) {
	b = &bot{
		store:            store,
		api:              api,
		genderController: &controllers.GenderController{},
		photoController:  &controllers.PhotoController{},
		logFile:          logFile,
		adminsList:       admins,
	}
	b.(*bot).setTimeLoggers()
	return b
}
