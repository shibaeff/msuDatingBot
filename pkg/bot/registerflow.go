package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	regBegin      = 0
	regName       = 1
	regGender     = 2
	regWantGender = 3
	regFaculty    = 4
	regAbout      = 5
	regPhoto      = 6
	regOver       = 7

	askName       = "Пожалуйста, введите свое имя"
	askGender     = "Пожалуйста, введите Ваш пол м/ж"
	askWantGender = "Кого ищем: м/ж?"
	askFaculty    = "С какого Вы факультета?"
	askAbout      = "Напишите немного о себе"
)

func registerStep(text string) (reply *tgbotapi.MessageConfig) {
	reply = &tgbotapi.MessageConfig{
		Text: text,
	}
	return
}

var UserQueue = make(map[int64]*User)

func registerFlow(message *tgbotapi.Message) (reply *tgbotapi.MessageConfig) {
	user, ok := UserQueue[message.Chat.ID]
	if !ok {
		user = &User{
			Id: message.Chat.ID,
		}
		UserQueue[message.Chat.ID] = user
	}
	switch RegisterStatus[message.Chat.ID] {
	case regBegin:
		RegisterStatus[message.Chat.ID] = regName
		return registerStep(askName)
	case regName:
		user.Name = message.Text
		log.Printf("Recorded username %s", user.Name)
		RegisterStatus[message.Chat.ID] = regGender
		reply = registerStep(askGender)
		return
	case regGender:
		RegisterStatus[message.Chat.ID] = regWantGender
		user.Gender = message.Text
		log.Printf("Recorded Gender %s", user.Gender)
		reply = registerStep(askWantGender)
	case regWantGender:
		user.WantGender = message.Text
		RegisterStatus[message.Chat.ID] = regFaculty
		log.Printf("Recorded want Gender %s", user.WantGender)
		reply = registerStep(askFaculty)
	case regFaculty:
		user.Faculty = message.Text
		RegisterStatus[message.Chat.ID] = regAbout
		log.Printf("Recorded Faculty %s", user.Faculty)
		reply = registerStep(askAbout)
	case regAbout:
		user.About = message.Text
		RegisterStatus[message.Chat.ID] = regPhoto
		log.Printf("Recorded about %s", user.About)
		reply = registerStep("Загрузите свое фото")
	case regPhoto:
		photos := *message.Photo
		photo := photos[0]
		user.PhotoLink = photo.FileID
		log.Printf("Recorded photo id %s", user.PhotoLink)
		RegisterStatus[message.Chat.ID] = regOver
		reply = registerStep("Регистрация окончена")
	}

	return
}
