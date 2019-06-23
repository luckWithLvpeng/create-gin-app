package controllers

import (
	"net/http"
	"strconv"

	"eme/models"
	"eme/pkg/code"

	"github.com/gin-gonic/gin"
)

// RoleGetByID role get by id
// @Summary 获取单个角色
// @Tags role
// @Security ApiKeyAuth
// @Param id path int true  "角色 id"
// @Success 200 {object} controllers.Response
// @Router /role/{id} [get]
func RoleGetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: code.InvalidParams,
			Msg:  code.GetMsg(code.InvalidParams) + err.Error(),
			Data: nil,
		})
		return
	}
	ecode, role, err := models.RoleGetByID(id)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: ecode,
			Msg:  code.GetMsg(ecode) + err.Error(),
			Data: nil,
		})
		return

	}
	c.JSON(http.StatusOK, Response{
		Code: ecode,
		Msg:  code.GetMsg(ecode),
		Data: role,
	})

}

// RoleGet role get
// @Summary 获取所有角色
// @Tags role
// @Security ApiKeyAuth
// @Success 200 {object} controllers.Response
// @Router /role [get]
func RoleGet(c *gin.Context) {
	ecode, data, err := models.RoleGet()
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: ecode,
			Msg:  code.GetMsg(ecode) + err.Error(),
			Data: nil,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code: ecode,
		Msg:  code.GetMsg(ecode),
		Data: data,
	})

}

// RoleAdd role add
// @Summary 添加角色
// @Tags role
// @Accept json
// @Security ApiKeyAuth
// @Param body body models.RoleForAdd true "role json"
// @Success 200 {object} controllers.Response
// @Router /role [post]
func RoleAdd(c *gin.Context) {
	var role models.Role
	err := c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: code.InvalidParams,
			Msg:  code.GetMsg(code.InvalidParams) + err.Error(),
			Data: nil,
		})
		return
	}
	ecode, id, err := models.RoleAdd(&role)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: ecode,
			Msg:  code.GetMsg(ecode) + err.Error(),
			Data: nil,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code: ecode,
		Msg:  code.GetMsg(ecode),
		Data: map[string]int64{
			"id": id,
		},
	})

}

// RoleUpdate role edit
// @Summary 编辑角色
// @Tags role
// @Accept json
// @Security ApiKeyAuth
// @Param body body models.RoleForAdd true "role json"
// @Param id path int true "role id"
// @Success 200 {object} controllers.Response
// @Router /role/{id} [put]
func RoleUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: code.InvalidParams,
			Msg:  code.GetMsg(code.InvalidParams) + err.Error(),
			Data: nil,
		})
		return
	}
	var role models.Role
	err = c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: code.InvalidParams,
			Msg:  code.GetMsg(code.InvalidParams) + err.Error(),
			Data: nil,
		})
		return
	}
	role.ID = int64(id)
	ecode, num, err := models.RoleUpdate(&role)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: ecode,
			Msg:  code.GetMsg(ecode) + err.Error(),
			Data: nil,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code: ecode,
		Msg:  code.GetMsg(ecode),
		Data: map[string]int64{
			"RowsAffected": num,
		},
	})

}

// RoleDel role delete
// @Summary 删除角色
// @Tags role
// @Accept x-www-form-urlencoded
// @Security ApiKeyAuth
// @Param id path int true "role id"
// @Success 200 {object} controllers.Response
// @Router /role/{id} [delete]
func RoleDel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusOK, Response{
			Code: code.InvalidParams,
			Msg:  code.GetMsg(code.InvalidParams) + err.Error(),
			Data: nil,
		})
		return
	}
	ecode, num, err := models.RoleDel(id)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: ecode,
			Msg:  code.GetMsg(ecode) + err.Error(),
			Data: nil,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code: ecode,
		Msg:  code.GetMsg(ecode),
		Data: map[string]int64{
			"RowsAffected": num,
		},
	})

}
