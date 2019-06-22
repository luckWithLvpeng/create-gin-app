package controllers

import (
	"net/http"
	"strconv"

	"eme/models"

	"github.com/gin-gonic/gin"
)

// UserGetByID user get by id
// @Summary 获取单个用户
// @Tags user
// @Accept  json
// @Param id path int true "user id"
// @Success 200 {object} controllers.Response
// @Router /user/{id} [get]
func UserGetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: InvalidParams,
			Msg:  GetMsg(InvalidParams) + err.Error(),
			Data: nil,
		})
		return
	}
	user, err := models.UserGetByID(id)
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
		Data: user,
	})

}

// UserGet user get
// @Summary 获取用户
// @Tags user
// @Accept  json
// @Param nowPage query int false "now page"
// @Param pageSize query int false "page size"
// @Success 200 {object} controllers.Response
// @Router /user [get]
func UserGet(c *gin.Context) {
	nowPage, err := strconv.Atoi(c.DefaultQuery("nowPage", "1"))
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "3"))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: InvalidParams,
			Msg:  GetMsg(InvalidParams) + err.Error(),
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
	count, users, err := models.UserGet(nowPage, pageSize)
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
		Data: map[string]interface{}{
			"users": users,
			"total": count,
		},
	})

}

// UserDeleteByID user delete by id
// @Summary 删除单个用户
// @Tags user
// @Accept  json
// @Param id path int true "user id"
// @Success 200 {object} controllers.Response
// @Router /user/{id} [delete]
func UserDeleteByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: InvalidParams,
			Msg:  GetMsg(InvalidParams) + err.Error(),
			Data: nil,
		})
		return
	}
	num, err := models.UserDelByID(int64(id))
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

// UserUpdate user update
// @Summary 更新用户
// @Tags user
// @Accept  json
// @Param id path int true "user id"
// @Param body body models.UserForUpdate ture "user for update"
// @Success 200 {object} controllers.Response
// @Router /user/{id} [put]
func UserUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: InvalidParams,
			Msg:  GetMsg(InvalidParams) + err.Error(),
			Data: nil,
		})
		return
	}
	var user models.UserForUpdate
	err = c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: InvalidParams,
			Msg:  GetMsg(InvalidParams) + err.Error(),
			Data: nil,
		})
		return
	}
	num, err := models.UserUpdate(id, &user)
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

// UserAdd user add
// @Summary 添加用户
// @Tags user
// @Accept  json
// @Param body body models.UserForAdd ture "user for add"
// @Success 200 {object} controllers.Response
// @Router /user [post]
func UserAdd(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: InvalidParams,
			Msg:  GetMsg(InvalidParams) + err.Error(),
			Data: nil,
		})
		return
	}
	id, err := models.UserAdd(&user)
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
