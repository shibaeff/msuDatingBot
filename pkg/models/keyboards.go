package models

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

var (
	genderMale = "М"
	genderFem  = ":"
	any        = "любой"
)
var (
	maleButton   = tgbotapi.InlineKeyboardButton{Text: "М", CallbackData: &genderMale}
	femaleBUtton = tgbotapi.InlineKeyboardButton{Text: "Ж", CallbackData: &genderFem}
	anyButton    = tgbotapi.InlineKeyboardButton{Text: any, CallbackData: &any}
)
var genderKeyboard = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(maleButton, femaleBUtton))
var wantGenderKeyboard = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(maleButton, femaleBUtton, anyButton))
