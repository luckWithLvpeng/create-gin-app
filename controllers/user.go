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
// @Tags user 用户
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
// @Tags user 用户
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

// UserRefreshToken user refresh token
// @Summary 用户获取新的token,新旧token会同时生效,旧的token 1分钟之后被销毁
// @Tags user 用户
// @Security ApiKeyAuth
// @Accept x-www-form-urlencoded
// @Param RefreshToken formData string true "refresh token"
// @Success 200 {object} controllers.Response
// @Router /user/refreshToken [post]
func UserRefreshToken(c *gin.Context) {
	Token := c.GetHeader("Authorization")
	RefreshToken := c.PostForm("RefreshToken")
	if Token == "" || RefreshToken == "" {
		c.JSON(http.StatusOK, Response{
			Code: code.InvalidParams,
			Msg:  code.GetMsg(code.InvalidParams),
			Data: nil,
		})
		return
	}
	ecode, token, err := models.UserRefreshToken(Token, RefreshToken)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: ecode,
			Msg:  code.GetMsg(ecode) + err.Error(),
			Data: token,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code: ecode,
		Msg:  code.GetMsg(ecode),
		Data: token,
	})

}

// UserLogout user logout ,add token to black list
// @Summary 用户退出, 销毁token 和 refreshToken
// @Tags user 用户
// @Accept x-www-form-urlencoded
// @Security ApiKeyAuth
// @Param RefreshToken formData string true "refresh token"
// @Success 200 {object} controllers.Response
// @Router /user/logout [post]
func UserLogout(c *gin.Context) {
	Token := c.GetHeader("Authorization")
	RefreshToken := c.PostForm("RefreshToken")
	if Token == "" || RefreshToken == "" {
		c.JSON(http.StatusOK, Response{
			Code: code.InvalidParams,
			Msg:  code.GetMsg(code.InvalidParams),
			Data: nil,
		})
		return
	}
	ecode, err := models.UserLogout(Token, RefreshToken)
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
		Data: nil,
	})
}

// UserLogin user login
// @Summary 用户登录
// @Tags user 用户
// @Accept x-www-form-urlencoded
// @Param Username formData string true "user name"
// @Param Password formData string true "user password"
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
// @Tags user 用户
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
// @Tags user 用户
// @Accept  json
// @Security ApiKeyAuth
// @Param id path int true "user id"
// @Param body body models.UserForUpdate true "user for update"
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
// @Tags user 用户
// @Accept  json
// @Security ApiKeyAuth
// @Param body body models.UserForAdd true "user for add"
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
