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
			id: message.Chat.ID,
		}
		UserQueue[message.Chat.ID] = user
	}
	switch RegisterStatus[message.Chat.ID] {
	case regBegin:
		RegisterStatus[message.Chat.ID] = regName
		return registerStep(askName)
	case regName:
		user.name = message.Text
		log.Printf("Recorded username %s", user.name)
		RegisterStatus[message.Chat.ID] = regGender
		reply = registerStep(askGender)
		return
	case regGender:
		RegisterStatus[message.Chat.ID] = regWantGender
		user.gender = message.Text
		log.Printf("Recorded gender %s", user.gender)
		reply = registerStep(askWantGender)
	case regWantGender:
		user.wantGender = message.Text
		RegisterStatus[message.Chat.ID] = regFaculty
		log.Printf("Recorded want gender %s", user.wantGender)
		reply = registerStep(askFaculty)
	case regFaculty:
		user.faculty = message.Text
		RegisterStatus[message.Chat.ID] = regAbout
		log.Printf("Recorded faculty %s", user.faculty)
		reply = registerStep(askAbout)
	case regAbout:
		user.about = message.Text
		RegisterStatus[message.Chat.ID] = regPhoto
		log.Printf("Recorded about %s", user.about)
		reply = registerStep("Загрузите свое фото")
	case regPhoto:
		photos := *message.Photo
		photo := photos[0]
		user.photoLink = photo.FileID
		log.Printf("Recorded photo id %s", user.photoLink)
		RegisterStatus[message.Chat.ID] = regOver
		reply = registerStep("Регистрация окончена")
	}

	return
}
