package main

import (
	"fmt"
	"github.com/pkg/errors"
	"homework/service"
	"log"
)

func main() {
	//直接模拟调用
	userlist, err := service.UserService.GetUserByAge(123)
	if err != nil {
		log.Fatalf("%v", errors.Cause(err))
	}
	fmt.Println(userlist)
	user, err := service.UserService.GetUserInfo(345)
	if err != nil {
		log.Fatalf("%v", err)
	} else {
		fmt.Println(user)
	}

}
