package models

import (
	"eme/middleware/auth"
	"eme/pkg/code"
	"encoding/json"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
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

// UserForLoginResponse 用户登录返回字段
type UserForLoginResponse struct {
	User         UserForGet
	Token        string
	Expires      int
	RefreshToken string
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
func UserGetByID(id int) (int, *UserForGet, error) {
	var user UserForGet
	DB := db
	DB = DB.Preload("Role").First(&user, id)
	if DB.Error == gorm.ErrRecordNotFound {
		log.Printf("can not find user: %v\n", DB.Error)
		return code.ErrorRecordNotFound, nil, DB.Error
	} else if DB.Error != nil {
		log.Printf("find user error: %v\n", DB.Error)
		return code.Error, nil, DB.Error
	}
	return code.Success, &user, nil
}

//UserDelByID 删除单个用户
func UserDelByID(id int64) (int, int64, error) {
	user := User{ID: id}
	DB := db
	DB = DB.Delete(&user)
	if DB.Error != nil {
		log.Printf("delete user error: %v\n", DB.Error)
		return code.Error, 0, DB.Error
	}
	auth.BlackList.Set(user.Username, user.Username, cache.DefaultExpiration)
	return code.Success, DB.RowsAffected, nil
}

//UserUpdate 更新用户
func UserUpdate(id int, u *UserForUpdate) (int, int64, error) {
	user := User{ID: int64(id)}
	DB := db
	DB = DB.First(&user)
	if user.Username == "" {
		return code.ErrorRecordNotFound, 0, errors.New("can not find the user")
	}
	DB = DB.Table("users").Where("id = ?", id).Updates(u)
	if DB.Error != nil {
		log.Printf("update user error: %v\n", DB.Error)
		return code.Error, 0, DB.Error
	}
	return code.Success, DB.RowsAffected, nil
}

//UserGet 获取用户
func UserGet(nowPage int, pageSize int) (int, int, []UserForGet, error) {
	var users []UserForGet
	var count int
	DB := db
	DB = DB.Table("users").Count(&count).Limit(pageSize).Offset((nowPage - 1) * pageSize).Preload("Role").Find(&users)
	if DB.Error != nil {
		log.Printf("find user error: %v\n", DB.Error)
		return code.Error, 0, nil, DB.Error
	}
	return code.Success, count, users, nil
}

// UserRefreshToken 用户刷新token
func UserRefreshToken(Token string, RefreshToken string) (int, string, error) {
	if _, found := auth.BlackList.Get(Token); found {
		return code.AuthInvalid, "", nil
	}
	_, err := auth.ParseToken(Token)
	if err != nil {
		return code.AuthInvalid, "", err
	}
	var user User
	DB := db
	DB.Where(&User{Key: RefreshToken}).Preload("Role").First(&user)
	if DB.Error == gorm.ErrRecordNotFound {
		return code.AuthInvalid, "", errors.New("refreshToken 失效")
	} else if DB.Error != nil {
		return code.AuthInvalid, "", DB.Error
	}
	newToken, err := auth.GenerateToken(user.Username, user.Role.RoleName)
	if err != nil {
		return code.Error, "", err
	}
	// 开启协程, 1分钟后销毁旧的token
	time.AfterFunc(time.Minute, func() {
		auth.BlackList.Set(Token, user.Username, cache.DefaultExpiration)
	})
	return code.Success, newToken, nil
}

// UserLogout 用户退出
func UserLogout(Token string, RefreshToken string) (int, error) {
	DB := db
	var user User
	DB.Where(&User{Key: RefreshToken}).First(&user)
	if DB.Error == gorm.ErrRecordNotFound {
		return code.AuthInvalid, errors.New("refreshToken 失效")
	} else if DB.Error != nil {
		return code.AuthInvalid, DB.Error
	}
	key := uuid.NewV4().String()
	DB = DB.Model(&user).Updates(User{Key: key})
	if DB.Error != nil {
		return code.Error, DB.Error
	}
	auth.BlackList.Set(Token, user.Username, cache.DefaultExpiration)
	auth.BlackList.Set(user.Username, user.Username, cache.DefaultExpiration)
	return code.Success, nil
}

// UserLogin 用户登录
func UserLogin(Username, Password string) (int, *UserForLoginResponse, error) {
	var user User
	DB := db
	DB.Where(&User{Username: Username}).Preload("Role").First(&user)
	if DB.Error == gorm.ErrRecordNotFound {
		return code.AuthInvalidUsernamePasssword, nil, DB.Error
	} else if DB.Error != nil {
		return code.Error, nil, DB.Error
	}
	// 比对密码是否匹配
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Password))
	if err != nil {
		return code.AuthInvalidUsernamePasssword, nil, err
	}

	key := uuid.NewV4().String()

	tmpLastLoginAt := user.LastLoginAt
	tx := DB.Begin()
	tx = tx.Model(&user).Updates(User{Key: key, Count: user.Count + 1, LastLoginAt: time.Now()})
	if tx.Error != nil {
		return code.Error, nil, tx.Error
	}
	token, err := auth.GenerateToken(Username, user.Role.RoleName)
	if err != nil {
		log.Printf("generate user token error: %v\n", tx.Error)
		tx.Rollback()
		return code.Error, nil, tx.Error
	}
	var res UserForLoginResponse
	userByte, _ := json.Marshal(user)
	json.Unmarshal(userByte, &res.User)
	res.Token = token
	res.RefreshToken = key
	res.User.LastLoginAt = tmpLastLoginAt
	res.Expires = auth.EffectiveDuration
	tx.Commit()
	// 从黑名单中把该用户去掉
	auth.BlackList.Delete(user.Username)
	return code.Success, &res, nil
}

//UserAdd 添加用户
func UserAdd(u *User) (int, *UserForLoginResponse, error) {
	var role Role
	DB := db
	DB = DB.Model(&u).Related(&role)
	if DB.Error != nil {
		log.Printf("can not find the role of user: %v\n", DB.Error)
		return code.Error, nil, DB.Error
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("hash user password error: %v\n", err)
		return code.Error, nil, err
	}
	u.Password = string(hashedPassword)
	u.LastLoginAt = time.Now()
	u.Role = role
	u.Key = uuid.NewV4().String()
	u.Count++
	// 开启事务
	tx := DB.Begin()
	tx = tx.Create(&u)
	if tx.Error != nil {
		log.Printf("creat user error: %v\n", tx.Error)
		tx.Rollback()
		return code.Error, nil, tx.Error
	}
	token, err := auth.GenerateToken(u.Username, u.Role.RoleName)
	if err != nil {
		log.Printf("generate user token error: %v\n", tx.Error)
		tx.Rollback()
		return code.Error, nil, tx.Error
	}
	var res UserForLoginResponse
	userByte, _ := json.Marshal(u)
	json.Unmarshal(userByte, &res.User)
	res.Token = token
	res.Expires = auth.EffectiveDuration
	res.RefreshToken = u.Key
	tx.Commit()
	return code.Success, &res, nil
}
