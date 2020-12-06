package bot

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"echoBot/pkg/models"
	"echoBot/pkg/store"
)

const (
	inlineMention = "[%s](tg://user?id=%d)"
)

func (b *bot) notifyUsers(message string) (list []*tgbotapi.MessageConfig, err error) {
	users, err := b.store.GetAllUsers()
	if err != nil {
		return
	}
	for _, user := range users {
		res := replyWithText(message)
		res.ChatID = user.Id
		list = append(list, res)
	}
	return
}

func (b *bot) ensureGender(u1, u2 *models.User) bool {
	return u1.Id != u2.Id &&
		u1.Gender == u2.WantGender && // strict pair
		u1.WantGender == u2.Gender || // another strict pair
		u1.WantGender == "любой" && u2.WantGender == u1.Gender ||
		u2.WantGender == "любой" && u1.WantGender == u2.Gender ||
		u1.WantGender == u2.WantGender && u1.WantGender == "любой"
}

func (b *bot) ensureAdmin(userName string) bool {
	for _, item := range b.adminsList {
		if item == userName {
			return true
		}
	}
	return false
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

func getUserLink(user *models.User) (raw string) {
	if user.UserName != "" {
		raw = fmt.Sprintf("@%s\n", user.UserName)
	} else {
		raw = fmt.Sprintf(inlineMention, user.Name, user.Id)
	}
	return
}
func (b *bot) prepareMatches(userId int64) (resp string, err error) {
	entry, err := b.store.GetMatchesRegistry().GetList(userId)
	if err != nil {
		return "Матчей нет", nil
	}
	if len(entry) == 0 {
		return "Матчей нет", nil
	}
	raw := []string{}
	for _, match := range entry {
		user, err := b.store.GetUser(match.Whome)
		if err != nil {
			continue
		}
		raw = append(raw, getUserLink(user))
	}
	resp = matchesList + strings.Join(raw, "\n")
	return
}

func replyWithText(text string) (ret *tgbotapi.MessageConfig) {
	ret = &tgbotapi.MessageConfig{
		Text: text,
	}
	ret.ReplyMarkup = menuKeyboard
	return
}

func replyWithPhoto(u *models.User, to int64) (ret *tgbotapi.PhotoConfig) {
	ret = &tgbotapi.PhotoConfig{
		BaseFile: tgbotapi.BaseFile{
			BaseChat: tgbotapi.BaseChat{
				ChatID: to,
			},
			UseExisting: true,
			FileID:      u.PhotoLink,
		},
		Caption: u.String(),
	}
	return
}

func find(slice []store.Entry, val int64) (int, bool) {
	for i, item := range slice {
		if item.Whome == val {
			return i, true
		}
	}
	return -1, false
}

func (b *bot) hide(user *models.User) {
	users, _ := b.store.GetAllUsers()
	for _, item := range users {
		if item.Id != user.Id {
			b.store.GetUnseenRegistry().DeleteItem(item.Id, user.Id)
			b.store.GetSeenRegistry().DeleteItem(item.Id, user.Id)
			b.store.GetMatchesRegistry().DeleteItem(item.Id, user.Id)
			b.store.GetLikesRegistry().DeleteItem(item.Id, user.Id)
		}
	}
}

func (b *bot) next(user *models.User) (reply interface{}) {
	if user.RegiStep < regOver {
		raw := replyWithText(notRegistered)
		raw.ChatID = user.Id
		reply = raw
		return
	}
	unseen, e := b.store.GetUnseen(user.Id)
	if len(unseen) == 0 || e != nil {
		reply = replyWithText("Вы просмотрели всех пользователей на данный момент")
		return reply
	}
	unseen_user, _ := b.store.GetUser(unseen[0].Whome)
	b.actionsLog.Printf("%d VIEWED %d\n", user.Id, unseen_user.Id)
	card := replyWithCard(unseen_user, user.Id)
	card.ParseMode = "html"
	reply = card
	return
}
