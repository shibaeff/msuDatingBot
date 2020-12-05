package main

import (
	"context"
	"io"
	"log"
	"os"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"echoBot/pkg/bot"
	"echoBot/pkg/store"
)

const (
	dbName              = "main"
	usersCollectionName = "users"
	likes               = "likes"
	unseen              = "unseen"
	matches             = "matches"
)

var (
	client, _ = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
)

func PrepareCollection(client *mongo.Client, name string) (conn *mongo.Collection) {
	conn = client.Database(dbName).Collection(name)
	return
}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	log.Println("hello")
	err := client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	users := PrepareCollection(client, usersCollectionName)
	likes := PrepareCollection(client, likes)
	unseen := PrepareCollection(client, unseen)
	matches := PrepareCollection(client, matches)
	seen := PrepareCollection(client, "seen")
	store := store.NewStore(users, []*mongo.Collection{likes, unseen, matches, seen})

	token, exists := os.LookupEnv("TELEGRAM_TOKEN")
	if !exists {
		log.Panic("Not telegram key specified!")
	}
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	readFile, err := os.OpenFile("log.txt", os.O_RDONLY, 0666)
	defer readFile.Close()
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	api.Debug = true

	log.Printf("Authorized on account %s", api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := api.GetUpdatesChan(u)
	admins, ok := os.LookupEnv("ADMINS")
	if !ok {
		log.Fatal("cannot load admins")
	}
	admins_list := strings.Split(admins, " ")
	Bot := bot.NewBot(store, api, readFile, admins_list)
	defer logFile.Close()
	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil { // ignore any non-Message Updates
			continue
		}

		var reply interface{}
		if update.Message != nil {
			reply, err = Bot.Reply(update.Message)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Println(update.CallbackQuery.Data)
			update.Message = &tgbotapi.Message{}
			update.Message.Text = update.CallbackQuery.Data
			update.Message.Chat = &tgbotapi.Chat{
				ID: int64(update.CallbackQuery.From.ID),
			}
			update.Message.From = &tgbotapi.User{}
			update.Message.From.UserName = update.CallbackQuery.From.UserName
			reply, err = Bot.Reply(update.Message)
		}

		switch v := reply.(type) {
		case *tgbotapi.MessageConfig:
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			//buttons := []tgbotapi.KeyboardButton{tgbotapi.KeyboardButton{Text: "Hello",},}
			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			//msg.ReplyToMessageID = update.Message.MessageID
			//msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons)
			v.ChatID = update.Message.Chat.ID
			v.ParseMode = "html"
			_, err = api.Send(v)
			if err != nil {
				log.Println(err)
			}
		case *tgbotapi.PhotoConfig:
			v.ParseMode = "html"
			tgbotapi.NewEditMessageText(1, 1, "")
			_, err := api.Send(v)
			if err != nil {
				log.Println(err)
			}
		case []*tgbotapi.MessageConfig:
			for _, item := range v {
				// item.ChatID = update.Message.Chat.ID
				_, err = api.Send(item)
				if err != nil {
					log.Println(err)
				}
			}
		case *bot.MessageToAdmin:
			config := v.MessageConfig
			for _, admin := range v.AdminsList {
				user, err := Bot.GetStore().GetAdmin(admin)
				if err == nil {
					config.ChatID = user.Id
					_, err = api.Send(config)
					if err != nil {
						log.Println(err)
					}
				}
			}
		case *bot.Match:
			_, err = api.Send(v.Msg1)
			if err != nil {
				log.Println(err)
			}
			_, err = api.Send(v.Msg2)
			if err != nil {
				log.Println(err)
			}
		}
		client.Ping(context.TODO(), nil)
	}
}
