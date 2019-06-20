package controllers

import (
	"net/http"
	"strconv"

	"eme/models"

	"github.com/gin-gonic/gin"
)

// UserGet user get
// @Summary 获取用户
// @Tags user
// @Accept  json
// @Param nowPage query int false "now page"
// @Param pageSize query int false "page size"
// @Success 200 {object} controllers.Response
// @Router /users [get]
func UserGet(c *gin.Context) {
	nowPage, err := strconv.Atoi(c.DefaultQuery("nowPage", "1"))
	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "1"))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			Code: InvalidParams,
			Msg:  GetMsg(InvalidParams) + err.Error(),
			Data: nil,
		})
		return
	}

	users, err := models.UserGet(nowPage, pageSize)
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
		},
	})

}

// UserAdd user add
// @Summary 添加用户
// @Tags user
// @Accept  json
// @Param body body models.UserFieldForAdd ture "user for add"
// @Success 200 {object} controllers.Response
// @Router /users [post]
func UserAdd(c *gin.Context) {
	var user models.UserFieldForAdd
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
