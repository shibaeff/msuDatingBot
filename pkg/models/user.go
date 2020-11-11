package models

import (
	"fmt"
)

const (
	stringify = "Имя %s\nФакультет%s\nПол%s\nПол собеседника%s\nО себе%s\n"
)

type User struct {
	Name       string
	Faculty    string
	Gender     string
	WantGender string
	About      string
	Id         int64
	PhotoLink  string
	RegiStep   int64
	UserName   string
}

func (u *User) String() string {
	return fmt.Sprintf(stringify, u.Name, u.Faculty, u.Gender, u.WantGender, u.About)
}
