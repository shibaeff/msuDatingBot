package store

import (
	"echoBot/pkg/bot"
)

const (
	initialSize = 100
)

type UserModel struct {
	User  *bot.User
	Seen  map[int64]bool
	Liked map[int64]bool
}

func NewUserModel(user *bot.User) *UserModel {
	return &UserModel{
		User:  user,
		Seen:  make(map[int64]bool),
		Liked: make(map[int64]bool),
	}
}
