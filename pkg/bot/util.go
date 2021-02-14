package bot

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"echoBot/pkg/models"
	"echoBot/pkg/store"
)

const (
	greetMsg1 = "Привет! ✨\nЭто бот знакомств МГУ. Работает аналогично Тиндеру 😉\n\nПожалуйста, перейди по этой [ссылке](%s) и ожидай подтверждения.\nПотом бот запросит имя, фоточку и пару слов о себе.\n\nПредложения и баги пишите в /feedback. "
	greetMsg2 = "Привет! ✨\nЭто бот знакомств МГУ. Работает аналогично Тиндеру 😉\nПотом бот запросит имя, фоточку и пару слов о себе.\nПредложения и баги пишите в /feedback.\nПоддержите нас донатом [здесь](https://vk.me/moneysend/cheptil)"
	linkStub  = "https://oauth.vk.com/authorize?client_id=7679100&scope=327682&&display=page&response_type=code&v=5.126&state=123456&redirect_uri=https://umsu.me/check?tg_id=%d"

	donateMsg = "Поддержите нас донатом [здесь](https://vk.me/moneysend/cheptil)"
)

func (b *bot) deleteUser(id int64) *tgbotapi.MessageConfig {
	b.store.DeleteUser(id)
	b.store.GetActions().DeleteEvents(store.Options{bson.E{"who", id}})
	b.store.GetActions().DeleteEvents(store.Options{bson.E{"whome", id}})
	reply := &tgbotapi.MessageConfig{}
	reply.Text = "Успешное удаление"
	reply.ChatID = id
	return reply
}

func (b *bot) notifyUsers(message string) (list []*tgbotapi.MessageConfig, err error) {
	users, err := b.store.GetAllUsers()
	if err != nil {
		return
	}
	for _, user := range users {
		res := user.ReplyWithText(message)
		b.api.Send(res)
	}
	return
}

func EnsureGender(u1, u2 *models.User) bool {
	return u1.Id != u2.Id &&
		(u1.Gender == u2.WantGender && // strict pair
			u1.WantGender == u2.Gender || // another strict pair
			u1.WantGender == "любой" && u2.WantGender == u1.Gender ||
			u2.WantGender == "любой" && u1.WantGender == u2.Gender ||
			u1.WantGender == u2.WantGender && u1.WantGender == "любой")
}

func (b *bot) grabLogs(offset int) (str string, err error) {
	var (
		part   []byte
		prefix bool
	)
	var txtlines []string
	reader := bufio.NewReader(b.logFile)
	buffer := bytes.NewBuffer(make([]byte, 1024))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			txtlines = append(txtlines, buffer.String())
			buffer.Reset()
		}
	}
	b.logFile.Seek(0, os.SEEK_SET)
	if len(txtlines) < offset {
		return "", errors.New("Неправильный оффсет")
	}
	err = nil
	str = strings.Join(txtlines[len(txtlines)-offset:], "\n")
	return
}

func (b *bot) parseLikee(message *tgbotapi.Message) (id int64, err error) {
	if message.ReplyToMessage == nil {
		return -1, errors.New("nothing to reply to")
	}
	text := message.ReplyToMessage.Text
	_, err = fmt.Scanf(text, &id)
	if err != nil {
		return -1, err
	}
	return
}

func (b *bot) prepareMatches(userId int64) (resp string, err error) {
	entry, err := b.store.GetActions().GetEvents(store.Options{
		bson.E{
			"who", userId,
		},
		bson.E{
			"event", store.EventMatch,
		},
	})
	if err != nil {
		return "Мэтчей нет", nil
	}
	if len(entry) == 0 {
		return "Мэтчей нет", nil
	}
	raw := []string{}
	for _, match := range entry {
		user, err := b.store.GetUser(match.Whome)
		if err != nil {
			continue
		}
		raw = append(raw, user.GetLink())
	}
	resp = matchesList + strings.Join(raw, "\n")
	return
}

func prepareHello(id int64) tgbotapi.MessageConfig {
	link := fmt.Sprintf(linkStub, id)
	msg := fmt.Sprintf(greetMsg1, link)
	// msg := fmt.Sprintf(greetMsg2)
	hello := tgbotapi.NewMessage(id, msg)
	hello.ParseMode = tgbotapi.ModeMarkdown
	return hello
}

func prepareDonate(id int64) tgbotapi.MessageConfig {
	msg := fmt.Sprintf(donateMsg)
	hello := tgbotapi.NewMessage(id, msg)
	hello.ParseMode = tgbotapi.ModeMarkdown
	return hello
}
