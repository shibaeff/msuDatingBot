package controllers

import (
	"errors"
)

var (
	GenderError  = errors.New("Invalid gender")
	PhotoError   = errors.New("Photo required")
	AboutError   = errors.New("About message too short")
	FacultyError = errors.New("No such faculty")
)

var (
	GenderErrorResp = "Введите правильный пол!"
	PhotoErrorResp  = "Отправьте валидное фото!"
	AboutErrorResp  = "Рассказ о себе должен быть длиннее 120 символов!"
	FacultyResp     = "Нет такого факультета!"
)
