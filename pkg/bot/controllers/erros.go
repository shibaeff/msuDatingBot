package controllers

import (
	"errors"
)

var (
	GenderError = errors.New("Invalid gender")
	PhotoError  = errors.New("Photo required")
	AboutError  = errors.New("About message too short")
)

var (
	GenderErrorResp = "Введите правильный пол!"
	PhotoErrorResp  = "Отправьте валидное фото!"
	AboutErrorResp  = "Рассказ о себе должен быть длиннее 120 символов!"
)
