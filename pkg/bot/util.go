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
	return u1.Id != u2.Id && u1.Gender == u2.WantGender && u1.WantGender == u2.Gender
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
		raw = append(raw, fmt.Sprintf("@%s\n", user.UserName))
	}
	resp = matchesList + strings.Join(raw, "")
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

func find(slice []int64, val int64) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
