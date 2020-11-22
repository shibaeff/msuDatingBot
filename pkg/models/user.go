package models

import (
	"fmt"
)

const (
	stringify = "*Имя* %s\n**Факультет** %s\n**Пол** %s\n**Пол собеседника** %s\n**О себе** %s\nUsername %s\nID %d"
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
	return fmt.Sprintf(stringify, u.Name, u.Faculty, u.Gender, u.WantGender, u.About, u.UserName, u.Id)
}
