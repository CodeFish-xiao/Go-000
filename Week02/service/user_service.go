package service

import "Week02/dao"

type IUserService interface {
	//获取y用户详细信息
	GetUserInfo(userId int64) (u dao.User, err error)
	//根据年龄查询列表
	GetUserByAge(age int) (u []dao.User, err error)
}
type UserServiceImpl struct {
}

var UserService IUserService

func (i UserServiceImpl) GetUserInfo(userId int64) (u dao.User, err error) {
	u.Id = userId
	err = u.GetUserById()
	return
}
func (i UserServiceImpl) GetUserByAge(age int) ([]dao.User, error) {
	return dao.GetUserByAge(age)
}
