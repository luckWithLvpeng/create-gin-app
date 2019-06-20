package models

// Role 角色
type Role struct {
	ID          int64 `gorm:"primary_key"`
	RoleName    string
	Description string
}

// RoleForAdd 添加角色
type RoleForAdd struct {
	RoleName    string
	Description string
}

func (Role) TableName() string {
	return "roles"
}

// RoleGetByID 获取单个角色
func RoleGetByID(id int) (*Role, error) {
	var role Role
	db := db.First(&role, id)
	if db.Error != nil {
		return nil, db.Error
	}
	return &role, err
}

// RoleGet 获取所有角色
func RoleGet() ([]*Role, error) {
	var roles []*Role
	db := db.Find(&roles)
	if db.Error != nil {
		return nil, db.Error
	}
	return roles, err
}

// RoleDel 删除角色
func RoleDel(id int) (int64, error) {
	role := Role{ID: int64(id)}
	db := db.Delete(&role)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, err
}

// RoleUpdate 编辑角色
func RoleUpdate(r *Role) (int64, error) {
	db := db.Save(&r)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil

}

// RoleAdd 添加角色
func RoleAdd(r *Role) (int64, error) {
	db := db.Create(&r)
	if db.Error != nil {
		return 0, db.Error
	}
	return r.ID, nil
}
