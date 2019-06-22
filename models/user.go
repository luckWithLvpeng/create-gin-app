package models

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// User 用户
type User struct {
	ID        int64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Role      Role
	RoleID    int64
	Email     string `gorm:"unique"`
	// jwt 通过对比key， 验证用户是否修改过密码
	Key         string
	Mobile      string `gorm:"type:varchar(30);unique"`
	LastLoginAt time.Time
	Disabled    bool `gorm:"type:tinyint"`
	// 登录次数
	Count int64
}

// UserForAdd 添加用户所需字段
type UserForAdd struct {
	Username string
	Password string
	RoleID   int64
	Email    string
	Mobile   string
}

// UserForUpdate 添加用户所需字段
type UserForUpdate struct {
	RoleID int64
	Email  string
	Mobile string
}

// UserForGet 获取用户返回字段
type UserForGet struct {
	ID          int64
	Username    string
	Role        Role
	RoleID      int64
	Disabled    bool
	Email       string
	Mobile      string
	LastLoginAt time.Time
	Count       int64
}

// TableName 指定用户表的名称
func (UserForUpdate) TableName() string {
	return "users"
}

// TableName 指定用户表的名称
func (UserForGet) TableName() string {
	return "users"
}

// TableName 指定用户表的名称
func (User) TableName() string {
	return "users"
}

//UserGetByID 获取单个用户
func UserGetByID(id int) (*UserForGet, error) {
	var user UserForGet
	DB := db
	DB = DB.Preload("Role").First(&user, id)
	if DB.Error == gorm.ErrRecordNotFound {
		log.Printf("find user error: %v\n", DB.Error)
		return nil, nil
	} else if DB.Error != nil {
		log.Printf("find user error: %v\n", DB.Error)
		return nil, DB.Error
	}
	return &user, nil
}

//UserDelByID 删除单个用户
func UserDelByID(id int64) (int64, error) {
	user := User{ID: id}
	DB := db
	DB = DB.Delete(&user)
	if DB.Error != nil {
		log.Printf("find user error: %v\n", DB.Error)
		return 0, DB.Error
	}
	return DB.RowsAffected, nil
}

//UserUpdate 更新用户
func UserUpdate(id int, u *UserForUpdate) (int64, error) {
	user := User{ID: int64(id)}
	DB := db
	DB = DB.First(&user)
	if user.Username == "" {
		return 0, errors.New("can not find the user")
	}
	DB = DB.Table("users").Where("id = ?", id).Updates(u)
	if DB.Error != nil {
		log.Printf("update user error: %v\n", DB.Error)
		return 0, DB.Error
	}
	return DB.RowsAffected, nil
}

//UserGet 获取用户
func UserGet(nowPage int, pageSize int) (int, []UserForGet, error) {
	var users []UserForGet
	var count int
	DB := db
	DB = DB.Table("users").Count(&count).Limit(pageSize).Offset((nowPage - 1) * pageSize).Preload("Role").Find(&users)
	if DB.Error != nil {
		log.Printf("find user error: %v\n", DB.Error)
		return 0, nil, DB.Error
	}
	return count, users, nil
}

//UserAdd 添加用户
func UserAdd(u *User) (int64, error) {
	var role Role
	DB := db
	DB = DB.Model(&u).Related(&role)
	if DB.Error != nil {
		log.Printf("can not find the role of user: %v\n", DB.Error)
		return 0, DB.Error
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("hash user password error: %v\n", err)
		return 0, err
	}
	u.Password = string(hashedPassword)
	u.LastLoginAt = time.Now()
	u.Role = role
	DB = db.Create(&u)
	if DB.Error != nil {
		log.Printf("creat user error: %v\n", DB.Error)
		return 0, DB.Error
	}
	return u.ID, nil
}
