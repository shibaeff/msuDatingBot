package controllers

import (
	"io/ioutil"
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

type AboutController struct {
}

func (a *AboutController) Verify(data interface{}) (string, error) {
	if len(data.(string)) < 120 {
		return AboutErrorResp, AboutError
	}
	return "", nil
}

type FacultyController struct {
	faculties map[string]bool
}

func NewFacultyController() Controller {
	f := &FacultyController{}
	f.faculties = make(map[string]bool)
	raw, err := ioutil.ReadFile("faculties.txt")
	if err != nil {
		panic(err)
	}
	all := string(raw)
	all = all[1 : len(all)-1]
	arrayOfNames := strings.ReplaceAll(all, "[", "")
	arrayOfNames = strings.ReplaceAll(arrayOfNames, "]", "")
	for _, item := range strings.Split(arrayOfNames, ",") {
		f.faculties[strings.ToLower(strings.Trim(item, " \n\""))] = true
	}
	return f
}

func (f *FacultyController) Verify(data interface{}) (string, error) {
	d := data.(string)
	_, ok := f.faculties[d]
	if !ok {
		return FacultyResp, FacultyError
	}
	return "", nil
}
