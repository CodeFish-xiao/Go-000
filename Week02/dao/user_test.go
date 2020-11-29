package dao

import (
	"fmt"
	"log"
	"testing"
)

func TestUser_GetUserById(t *testing.T) {
	user := User{Id: 4}
	err := user.GetUserById()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)
}
