package auth

import (
	"eme/pkg/code"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Auth 验证用户权限的中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ecode := code.Success
		token := c.GetHeader("Authorization")
		if token == "" {
			ecode = code.AuthNeedHeaderAuthorization
		} else if _, found := BlackList.Get(token); found {
			// token 已经失效
			ecode = code.AuthInvalid
		} else {
			claims, err := ParseToken(token)
			if err != nil {
				ecode = code.AuthParseToken
			} else if _, found := BlackList.Get(claims.Username); found {
				ecode = code.AuthInvalid
			} else if time.Now().Unix() > claims.ExpiresAt {
				ecode = code.AuthTokenTimeout
			}
		}

		if ecode != code.Success {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Code": ecode,
				"Msg":  code.GetMsg(ecode),
				"Data": nil,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
