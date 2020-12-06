package service

import "sync"

var serviceOnce sync.Once

func init() {
	serviceOnce.Do(func() {
		UserService = &UserServiceImpl{}
	})
}
