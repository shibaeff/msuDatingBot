package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"echoBot/pkg/models"
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

func (b *bot) updateRegStatus(id, status int64) error {
	return b.store.UpdUserField(id, "registep", status)
}

func (b *bot) registerFlow(user *models.User, message *tgbotapi.Message) (reply *tgbotapi.MessageConfig) {
	switch user.RegiStep {
	case waiting:
		if err := b.updateRegStatus(user.Id, regName); err != nil {
			log.Fatal(err)
		}
		return registerStep(askName)
	case regBegin:
		// RegisterStatus[message.Chat.ID] = regName
		if err := b.updateRegStatus(user.Id, regName); err != nil {
			log.Fatal(err)
		}
		return registerStep(askName)
	case regName:
		user.Name = message.Text
		log.Printf("Recorded username %s", user.Name)
		// RegisterStatus[message.Chat.ID] = regGender
		if err := b.updateRegStatus(user.Id, regGender); err != nil {
			log.Fatal(err)
		}
		reply = registerStep(askGender)
		return
	case regGender:
		// RegisterStatus[message.Chat.ID] = regWantGender
		if err := b.updateRegStatus(user.Id, regWantGender); err != nil {
			log.Fatal(err)
		}
		user.Gender = message.Text
		log.Printf("Recorded Gender %s", user.Gender)
		reply = registerStep(askWantGender)
	case regWantGender:
		user.WantGender = message.Text
		// RegisterStatus[message.Chat.ID] = regFaculty
		if err := b.updateRegStatus(user.Id, regFaculty); err != nil {
			log.Fatal(err)
		}
		log.Printf("Recorded want Gender %s", user.WantGender)
		reply = registerStep(askFaculty)
	case regFaculty:
		user.Faculty = message.Text
		// RegisterStatus[message.Chat.ID] = regAbout
		if err := b.updateRegStatus(user.Id, regAbout); err != nil {
			log.Fatal(err)
		}
		log.Printf("Recorded Faculty %s", user.Faculty)
		reply = registerStep(askAbout)
	case regAbout:
		user.About = message.Text
		// RegisterStatus[message.Chat.ID] = regPhoto
		if err := b.updateRegStatus(user.Id, regPhoto); err != nil {
			log.Fatal(err)
		}
		log.Printf("Recorded about %s", user.About)
		reply = registerStep("Загрузите свое фото")
	case regPhoto:
		photos := *message.Photo
		photo := photos[0]
		user.PhotoLink = photo.FileID
		log.Printf("Recorded photo id %s", user.PhotoLink)
		// RegisterStatus[message.Chat.ID] = regOver
		if err := b.updateRegStatus(user.Id, regOver); err != nil {
			log.Fatal(err)
		}
		reply = registerStep("Регистрация окончена")
	}

	return
}
