package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"echoBot/pkg/bot"
	"echoBot/pkg/store"
)

const (
	dbName                = "main"
	usersCollectionName   = "users"
	actionsCollectionName = "actions"
)

var (
	client, _ = mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
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

func switchReply(api *tgbotapi.BotAPI, reply interface{}) {
	switch v := reply.(type) {
	case *tgbotapi.MessageConfig:
		if _, err := api.Send(v); err != nil {
			log.Println(err)
		}
	case *tgbotapi.PhotoConfig:
		api.Send(v)
	case *tgbotapi.DocumentConfig:
		api.Send(v)
	}
}

func main() {
	//go func() {
	//	err := http.ListenAndServe(":3000", nil)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}()
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
	actions := PrepareCollection(client, actionsCollectionName)
	store := store.NewStore(users, actions)

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

	// updates, err := api.GetUpdatesChan(u)
	admins, ok := os.LookupEnv("ADMINS")
	if !ok {
		log.Fatal("cannot load admins")
	}
	admins_list := strings.Split(admins, " ")
	Bot := bot.NewBot(store, api, readFile, admins_list)
	//api.SetWebhook(tgbotapi.NewWebhookWithCert(os.Getenv("WEBHOOK", )))
	defer logFile.Close()
	updates := api.ListenForWebhook("/")
	go http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)
	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil { // ignore any non-Message Updates
			continue
		}

		var reply interface{}
		if update.Message != nil {
			ctx, _ := context.WithTimeout(context.Background(), time.Second)
			go func() {
				reply, err = Bot.ReplyMessage(ctx, update.Message)
				if err != nil {
					log.Fatal(err)
				}
				switchReply(api, reply)
			}()
		} else {
			ctx, _ := context.WithTimeout(context.Background(), time.Second)
			go func() {
				reply, _ = Bot.HandleCallbackQuery(ctx, update.CallbackQuery)
				switchReply(api, reply)
				if err != nil {
					log.Fatal(err)
				}
			}()
		}
	}
}
