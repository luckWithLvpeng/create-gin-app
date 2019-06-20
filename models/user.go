package models

import (
	"time"
)

// User 用户
type User struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name string
	Role Role
	RoleID int64
	Email string
	LastLoginAt time.Time
	Count  int64  // 登录次数
}

// UserFieldForAdd 添加用户所需字段
type UserFieldForAdd struct {
	Name string `db:"name"`
	Role int64  `db:"role"`
}
func (User) TableName() string {
	return "users"
}
//UserGet 获取用户
func UserGet(nowPage int, pageSize int) ([]User, error) {
	//tx := DB.MustBegin()

	return nil, nil

}

//UserAdd 添加用户
func UserAdd(u *UserFieldForAdd) (int64, error) {

	return 0, nil

}
