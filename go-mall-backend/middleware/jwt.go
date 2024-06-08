package middleware

import (
	"github.com/gin-gonic/gin"
	"mall/pkg/e"
	"mall/pkg/utils"
	"time"
)

// JWT
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int         // Status code to be returned
		var data interface{} // Data to be returned
		code = 200
		token := c.GetHeader("Authorization") // Get token from Authorization header
		if token == "" {
			code = 404 // No token found
		} else {
			claims, err := util.ParseToken(token) // Parse token to get claims
			if err != nil {
				code = e.ErrorAuthCheckTokenFail // Token parsing error
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ErrorAuthCheckTokenTimeout // Token has expired
			}
		}
		// If there is an error, respond with the error code and message and abort the request
		if code != e.SUCCESS {
			c.JSON(200, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   data,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// JWTAdmin is middleware for validating JWT tokens and checking for admin privileges
func JWTAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}
		token := c.GetHeader("Authorization")
		if token == "" {
			code = e.InvalidParams
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ErrorAuthCheckTokenFail
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ErrorAuthCheckTokenTimeout
			} else if claims.Authority == 0 {
				code = e.ErrorAuthInsufficientAuthority
			}
		}
		if code != e.SUCCESS {
			c.JSON(200, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   data,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
