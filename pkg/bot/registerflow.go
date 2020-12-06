package bot

import (
	"log"
	"strings"

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
	askWantGender = "Кого ищем: м/ж/любой?"
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
		if err := b.store.UpdUserField(user.Id, "name", user.Name); err != nil {
			log.Fatal(err)
		}
		log.Printf("Recorded username %s", user.Name)
		// RegisterStatus[message.Chat.ID] = regGender
		if err := b.updateRegStatus(user.Id, regGender); err != nil {
			log.Fatal(err)
		}
		reply = registerStep(askGender)
		return
	case regGender:
		// RegisterStatus[message.Chat.ID] = regWantGender
		resp, err := b.genderController.Verify(message.Text)
		if err != nil {
			reply = replyWithText(resp)
			return
		}
		if err := b.updateRegStatus(user.Id, regWantGender); err != nil {
			log.Fatal(err)
		}
		user.Gender = message.Text
		if err := b.store.UpdUserField(user.Id, "gender", strings.ToLower(user.Gender)); err != nil {
			log.Fatal(err)
		}
		log.Printf("Recorded Gender %s", user.Gender)

		reply = registerStep(askWantGender)
	case regWantGender:
		resp, err := b.genderController.Verify(message.Text)
		if err != nil {
			reply = replyWithText(resp)
			return
		}
		user.WantGender = message.Text
		// RegisterStatus[message.Chat.ID] = regFaculty
		if err := b.updateRegStatus(user.Id, regFaculty); err != nil {
			log.Fatal(err)
		}
		if err := b.store.UpdUserField(user.Id, "wantgender", strings.ToLower(user.WantGender)); err != nil {
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
		if err := b.store.UpdUserField(user.Id, "faculty", user.Faculty); err != nil {
			log.Fatal(err)
		}
		log.Printf("Recorded Faculty %s", user.Faculty)
		reply = registerStep(askAbout)
	case regAbout:
		user.About = message.Text
		resp, err := b.aboutController.Verify(user.About)
		if err != nil {
			return replyWithText(resp)
		}
		// RegisterStatus[message.Chat.ID] = regPhoto
		if err := b.updateRegStatus(user.Id, regPhoto); err != nil {
			log.Fatal(err)
		}
		if err := b.store.UpdUserField(user.Id, "about", user.About); err != nil {
			log.Fatal(err)
		}
		log.Printf("Recorded about %s", user.About)
		reply = registerStep("Загрузите свое фото")
	case regPhoto:
		resp, err := b.photoController.Verify(message.Photo == nil)
		if err != nil {
			reply = replyWithText(resp)
			return
		}
		photos := *message.Photo
		photo := photos[0]
		user.PhotoLink = photo.FileID
		if err := b.store.UpdUserField(user.Id, "photolink", user.PhotoLink); err != nil {
			log.Fatal(err)
		}
		log.Printf("Recorded photo id %s", user.PhotoLink)
		// RegisterStatus[message.Chat.ID] = regOver
		if err := b.updateRegStatus(user.Id, regOver); err != nil {
			log.Fatal(err)
		}
		user.RegiStep = regOver
		reply = registerStep("Регистрация окончена\n" + user.String())
	}

	return
}
