package controllers

import (
	"net/http"
	"strconv"

	"eme/models"

	"github.com/gin-gonic/gin"
)

// RoleGetByID role get by id
// @Summary 获取单个角色
// @Tags role
// @Param id path int true  "角色 id"
// @Success 200 {object} controllers.Response
// @Router /role/{id} [get]
func RoleGetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: InvalidParams,
			Msg:  GetMsg(InvalidParams) + err.Error(),
			Data: nil,
		})
		return
	}
	role, err := models.RoleGetByID(id)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: Error,
			Msg:  GetMsg(Error) + err.Error(),
			Data: nil,
		})
		return

	}
	c.JSON(http.StatusOK, Response{
		Code: Success,
		Msg:  GetMsg(Success),
		Data: role,
	})

}

// RoleGet role get
// @Summary 获取所有角色
// @Tags role
// @Success 200 {object} controllers.Response
// @Router /role [get]
func RoleGet(c *gin.Context) {
	data, err := models.RoleGet()
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: Error,
			Msg:  GetMsg(Error) + err.Error(),
			Data: nil,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code: Success,
		Msg:  GetMsg(Success),
		Data: data,
	})

}

// RoleAdd role add
// @Summary 添加角色
// @Tags role
// @Accept json
// @Param body body models.RoleForAdd true "role json"
// @Success 200 {object} controllers.Response
// @Router /role [post]
func RoleAdd(c *gin.Context) {
	var role models.Role
	err := c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: InvalidParams,
			Msg:  GetMsg(InvalidParams) + err.Error(),
			Data: nil,
		})
		return
	}
	id, err := models.RoleAdd(&role)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: Error,
			Msg:  GetMsg(Error) + err.Error(),
			Data: nil,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code: Success,
		Msg:  GetMsg(Success),
		Data: map[string]int64{
			"id": id,
		},
	})

}

// RoleUpdate role edit
// @Summary 编辑角色
// @Tags role
// @Accept json
// @Param body body models.RoleForAdd true "role json"
// @Param id path int true "role id"
// @Success 200 {object} controllers.Response
// @Router /role/{id} [put]
func RoleUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: InvalidParams,
			Msg:  GetMsg(InvalidParams) + err.Error(),
			Data: nil,
		})
		return
	}
	var role models.Role
	err = c.ShouldBindJSON(&role)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: InvalidParams,
			Msg:  GetMsg(InvalidParams) + err.Error(),
			Data: nil,
		})
		return
	}
	role.ID = int64(id)
	num, err := models.RoleUpdate(&role)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: Error,
			Msg:  GetMsg(Error) + err.Error(),
			Data: nil,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code: Success,
		Msg:  GetMsg(Success),
		Data: map[string]int64{
			"RowsAffected": num,
		},
	})

}

// RoleDel role delete
// @Summary 删除角色
// @Tags role
// @Accept x-www-form-urlencoded
// @Param id path int true "role id"
// @Success 200 {object} controllers.Response
// @Router /role/{id} [delete]
func RoleDel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusOK, Response{
			Code: InvalidParams,
			Msg:  GetMsg(InvalidParams) + err.Error(),
			Data: nil,
		})
		return
	}
	num, err := models.RoleDel(id)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: Error,
			Msg:  GetMsg(Error) + err.Error(),
			Data: nil,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code: Success,
		Msg:  GetMsg(Success),
		Data: map[string]int64{
			"RowsAffected": num,
		},
	})

}
