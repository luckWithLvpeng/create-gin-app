package models

import (
	"create-gin-app/pkg/code"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Role 角色
type Role struct {
	ID          int64  `gorm:"primary_key"`
	RoleName    string `gorm:"unique;not null"`
	Description string
}

// RoleForAdd 添加角色
type RoleForAdd struct {
	RoleName    string
	Description string
}

// TableName 获取表名
func (Role) TableName() string {
	return "roles"
}

// RoleGetByID 获取单个角色
func RoleGetByID(id int) (int, *Role, error) {
	var role Role
	DB := db
	DB = DB.First(&role, id)
	if DB.Error == gorm.ErrRecordNotFound {
		log.Printf("can not find role: %v\n", DB.Error)
		return code.ErrorRecordNotFound, nil, DB.Error
	} else if DB.Error != nil {
		log.Printf("find role error: %v\n", DB.Error)
		return code.Error, nil, DB.Error
	}
	return code.Success, &role, nil
}

// RoleGet 获取所有角色
func RoleGet() (int, []*Role, error) {
	var roles []*Role
	DB := db
	DB = DB.Find(&roles)
	if DB.Error != nil {
		return code.Error, nil, DB.Error
	}
	return code.Success, roles, nil
}

// RoleDel 删除角色
func RoleDel(id int) (int, int64, error) {
	role := Role{ID: int64(id)}
	DB := db
	DB = DB.Delete(&role)
	if DB.Error != nil {
		return code.Error, 0, DB.Error
	}
	return code.Success, DB.RowsAffected, nil
}

// RoleUpdate 编辑角色
func RoleUpdate(r *Role) (int, int64, error) {
	role := Role{ID: r.ID}
	DB := db
	DB = DB.First(&role)
	if role.RoleName == "" {
		return code.Error, 0, errors.New("can not find the role")
	}
	DB = DB.Save(&r)
	if DB.Error != nil {
		return code.Error, 0, DB.Error
	}
	return code.Success, DB.RowsAffected, nil

}

// RoleAdd 添加角色
func RoleAdd(r *Role) (int, int64, error) {
	var role Role
	db.Where("role_name = ?", r.RoleName).First(&role)
	if role.ID > 0 {
		return code.Error, 0, errors.New("role already exist")
	}
	db := db.Create(&r)
	if db.Error != nil {
		return code.Error, 0, db.Error
	}
	return code.Success, r.ID, nil
}
