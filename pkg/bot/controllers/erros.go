package controllers

import (
	"errors"
)

var (
	GenderError = errors.New("Invalid gender")
	PhotoError  = errors.New("Photo required")
)

var (
	GenderErrorResp = "Введите правильный пол!"
	PhotoErrorResp  = "Отправьте валидное фото!"
)
