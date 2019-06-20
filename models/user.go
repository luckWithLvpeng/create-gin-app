package models

import (
	"log"
	"time"
)

// User 用户
type User struct {
	ID          int64 `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Username    string
	Password    string
	Role        Role
	RoleID      int64
	Email       string
	Token       string
	mobile      string `gorm:"type:varchar(30)"`
	LastLoginAt time.Time
	Count       int64 // 登录次数
}

// UserForAdd 添加用户所需字段
type UserForAdd struct {
	Username string
	Password string
	RoleID   int64
	Email    string
	mobile   string
}

func (User) TableName() string {
	return "users"
}

//UserGet 获取用户
func UserGet(nowPage int, pageSize int) ([]User, error) {
	//tx := DB.MustBegin()
	var users []User
	db := db.Offset((nowPage-1) * pageSize).Limit(pageSize).Find(&users)
	if db.Error != nil {
		log.Printf("find user error: %v\n", err)
		return nil, db.Error
	}
	return users, nil
}

//UserAdd 添加用户
func UserAdd(u *User) (int64, error) {
	var role Role
	u.LastLoginAt = time.Now()
	db := db.Model(&u).Related(&role)
	if db.Error != nil {
		log.Printf("can not find the role of user: %v\n",  db.Error)
		return 0, db.Error
	}
	db = db.Create(&u)
	if db.Error != nil {
		log.Printf("creat user error: %v\n", err)
		return 0, db.Error
	}
	return u.ID, nil
}
