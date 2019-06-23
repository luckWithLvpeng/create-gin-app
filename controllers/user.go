package controllers

import (
	"net/http"
	"strconv"

	"eme/models"
	"eme/pkg/code"

	"github.com/gin-gonic/gin"
)

// UserGetByID user get by id
// @Summary 获取单个用户
// @Tags user
// @Accept  json
// @Security ApiKeyAuth
// @Param id path int true "user id"
// @Success 200 {object} controllers.Response
// @Router /user/{id} [get]
func UserGetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: code.InvalidParams,
			Msg:  code.GetMsg(code.InvalidParams) + err.Error(),
			Data: nil,
		})
		return
	}
	ecode, user, err := models.UserGetByID(id)
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
		Data: user,
	})

}

// UserGet user get
// @Summary 获取用户
// @Tags user
// @Accept  json
// @Security ApiKeyAuth
// @Param nowPage query int false "now page"
// @Param pageSize query int false "page size"
// @Success 200 {object} controllers.Response
// @Router /user [get]
func UserGet(c *gin.Context) {
	nowPage, err := strconv.Atoi(c.DefaultQuery("nowPage", "1"))
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "3"))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: code.InvalidParams,
			Msg:  code.GetMsg(code.InvalidParams) + err.Error(),
			Data: nil,
		})
		return
	}
	if nowPage <= 0 {
		nowPage = 1
	}
	if pageSize <= 0 {
		pageSize = 3
	}
	ecode, count, users, err := models.UserGet(nowPage, pageSize)
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
		Data: map[string]interface{}{
			"users": users,
			"total": count,
		},
	})

}

// UserLogout user logout
// @Summary 用户退出，销毁token
// @Tags user
// @Accept x-www-form-urlencoded
// @Security ApiKeyAuth
// @Success 200 {object} controllers.Response
// @Router /user/logout [post]
func UserLogout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	ecode := models.UserLogout(token)
	c.JSON(http.StatusOK, Response{
		Code: ecode,
		Msg:  code.GetMsg(ecode),
		Data: nil,
	})

}

// UserLogin user login
// @Summary 用户登录
// @Tags user
// @Accept x-www-form-urlencoded
// @Param Username formData string ture "user name"
// @Param Password formData string ture "user password"
// @Success 200 {object} controllers.Response
// @Router /user/login [post]
func UserLogin(c *gin.Context) {
	Username := c.PostForm("Username")
	Password := c.PostForm("Password")
	if Username == "" || Password == "" {
		c.JSON(http.StatusOK, Response{
			Code: code.InvalidParams,
			Msg:  code.GetMsg(code.InvalidParams),
			Data: nil,
		})
		return
	}
	ecode, data, err := models.UserLogin(Username, Password)
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

// UserDeleteByID user delete by id
// @Summary 删除单个用户
// @Tags user
// @Accept  json
// @Security ApiKeyAuth
// @Param id path int true "user id"
// @Success 200 {object} controllers.Response
// @Router /user/{id} [delete]
func UserDeleteByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: code.InvalidParams,
			Msg:  code.GetMsg(code.InvalidParams) + err.Error(),
			Data: nil,
		})
		return
	}
	ecode, num, err := models.UserDelByID(int64(id))
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

// UserUpdate user update
// @Summary 更新用户
// @Tags user
// @Accept  json
// @Security ApiKeyAuth
// @Param id path int true "user id"
// @Param body body models.UserForUpdate ture "user for update"
// @Success 200 {object} controllers.Response
// @Router /user/{id} [put]
func UserUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: code.InvalidParams,
			Msg:  code.GetMsg(code.InvalidParams) + err.Error(),
			Data: nil,
		})
		return
	}
	var user models.UserForUpdate
	err = c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: code.InvalidParams,
			Msg:  code.GetMsg(code.InvalidParams) + err.Error(),
			Data: nil,
		})
		return
	}
	ecode, num, err := models.UserUpdate(id, &user)
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

// UserAdd user add
// @Summary 添加用户
// @Tags user
// @Accept  json
// @Security ApiKeyAuth
// @Param body body models.UserForAdd ture "user for add"
// @Success 200 {object} controllers.Response
// @Router /user [post]
func UserAdd(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: code.InvalidParams,
			Msg:  code.GetMsg(code.InvalidParams) + err.Error(),
			Data: nil,
		})
		return
	}
	ecode, data, err := models.UserAdd(&user)
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
