package models

import (
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
func RoleGetByID(id int) (*Role, error) {
	var role Role
	DB := db
	DB = DB.First(&role, id)
	if DB.Error != nil {
		return nil, DB.Error
	}
	return &role, nil
}

// RoleGet 获取所有角色
func RoleGet() ([]*Role, error) {
	var roles []*Role
	DB := db
	DB = DB.Find(&roles)
	if DB.Error != nil {
		return nil, DB.Error
	}
	return roles, nil
}

// RoleDel 删除角色
func RoleDel(id int) (int64, error) {
	role := Role{ID: int64(id)}
	DB := db
	DB = DB.Delete(&role)
	if DB.Error != nil {
		return 0, DB.Error
	}
	return DB.RowsAffected, nil
}

// RoleUpdate 编辑角色
func RoleUpdate(r *Role) (int64, error) {
	role := Role{ID: r.ID}
	DB := db
	DB = DB.First(&role)
	if role.RoleName == "" {
		return 0, errors.New("can not find the role")
	}
	DB = DB.Save(&r)
	if DB.Error != nil {
		return 0, DB.Error
	}
	return DB.RowsAffected, nil

}

// RoleAdd 添加角色
func RoleAdd(r *Role) (int64, error) {
	var role Role
	db.Where("role_name = ?", r.RoleName).First(&role)
	if role.ID > 0 {
		return 0, errors.New("role already exist")
	}
	db := db.Create(&r)
	if db.Error != nil {
		return 0, db.Error
	}
	return r.ID, nil
}
