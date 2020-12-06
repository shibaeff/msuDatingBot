package controllers

import (
	"strings"
)

type GenderController struct {
}

func (g *GenderController) Verify(data interface{}) (s string, err error) {
	switch v := data.(type) {
	case string:
		sl := strings.ToLower(v)
		if sl != "м" && sl != "ж" && sl != "любой" {
			return GenderErrorResp, GenderError
		} else {
			return "", nil
		}
	default:
		return GenderErrorResp, GenderError
	}
}

type PhotoController struct {
}

func (p *PhotoController) Verify(data interface{}) (string, error) {
	if data.(bool) {
		return PhotoErrorResp, PhotoError
	}
	return "", nil
}
