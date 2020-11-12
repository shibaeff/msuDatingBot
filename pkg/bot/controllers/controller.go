package controllers

type Controller interface {
	Verify(data interface{}) (string, error)
}
