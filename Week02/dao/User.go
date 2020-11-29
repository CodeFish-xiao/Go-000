package dao

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var (
	db  *sql.DB
	err error
)

type User struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Age  string `json:"age"`
}

func init() {
	db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/go_00?charset=utf8")
	if err != nil {
		panic(err)
	}

}
func (u *User) GetUserById() error {
	//查询数据，指定字段名，返回sql.Rows结果集
	err = db.QueryRow("select id,name,age from user where id =? ", u.Id).Scan(&u.Id, &u.Name, &u.Age)
	if err != nil {
		return errors.Wrap(err, "数据库没有数据！")
	}
	return nil
}

//根据年龄查询
func GetUserByAge(age int) ([]User, error) {
	//查询数据，指定字段名，返回sql.Rows结果集
	rows, err := db.Query("select id,name,age from user where age =? ", age)
	var userList []User
	if err != nil {
		return nil, errors.Wrapf(err, "数据库错误")
	}
	//没数据就返回空数组
	for rows.Next() {
		u := User{}
		_ = rows.Scan(&u.Id, &u.Name, &u.Age)
		userList = append(userList, u)
	}
	return userList, nil
}

func Getdb() *sql.DB {
	return db
}
