package models

import (
	"fmt"
	"testing"
)

func TestUser_String(t *testing.T) {
	user := &User{Name: "Pasha"}
	fmt.Println(user)
}
