package models

import (
	"fmt"
)

const (
	stringify = "<strong>Имя:</strong> %s\n" +
		"<strong>Факультет:</strong> %s\n" +
		"<strong>Пол:</strong> %s\n" +
		"<strong>Пол собеседника:</strong> %s\n" +
		"<strong>О себе:</strong> %s\n"
	//"<strong>Username:</strong> %s\n"
	//"<strong>ID:</strong> %d"
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
