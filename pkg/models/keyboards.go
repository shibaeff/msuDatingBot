package models

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

var (
	genderMale = "male"
	genderFem  = "female"
)
var (
	maleButton   = tgbotapi.InlineKeyboardButton{Text: "М", CallbackData: &genderMale}
	femaleBUtton = tgbotapi.InlineKeyboardButton{Text: "Ж", CallbackData: &genderFem}
)
var genderKeyboard = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(maleButton, femaleBUtton))
