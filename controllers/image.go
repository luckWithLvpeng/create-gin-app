package controllers

import (
	"eme/pkg/code"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ImageClassify image classify
// @Summary 图像分类，识别通用物体
// @Description 输入一张图片，输出图片中的多个通用物体和场景, 给图片分类
// @Tags image 图像
// @Security ApiKeyAuth
// @Accept x-www-form-urlencoded
// @Param image formData string true  "图片的base64编码,去掉编码头,图片大小不超过4M"
// @Success 200 {object} controllers.Response
// @Router /image/classify [post]
func ImageClassify(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: code.Success,
		Msg:  code.GetMsg(code.Success),
		Data: nil,
	})

}

// ImageFaceDetect image face detect
// @Summary 人脸检测
// @Description 获取人脸的位置, 人脸相关的属性，如 性别，年龄等，人脸的质量信息，相关因素如亮度，遮挡，模糊，完整度， 置信度等；人脸关键点信息等
// @Tags image 图像
// @Security ApiKeyAuth
// @Accept x-www-form-urlencoded
// @Param image formData string true  "图片的base64编码,去掉编码头"
// @Success 200 {object} controllers.Response
// @Router /image/face/detect [post]
func ImageFaceDetect(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: code.Success,
		Msg:  code.GetMsg(code.Success),
		Data: nil,
	})

}
