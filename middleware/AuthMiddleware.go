package middleware

import (
	"GinProject12/common"
	"GinProject12/databases"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "权限不足"})
			c.Abort()
			return
		}

		tokenString = tokenString[7:]

		id, err := common.ParseJWT(tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "权限不足"})
			c.Abort()
			return
		}

		user := sdb.FindUserByEvery("username", id)

		admin := sdb.FindAdminByEvery("username", id)
		if admin == nil && user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "权限不足"})
			c.Abort()
			return
		} else if admin == nil {
			c.Set("user", *user)
		} else {
			c.Set("user", *admin)
		}
		c.Next()

	}
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "权限不足"})
			c.Abort()
			return
		}

		tokenString = tokenString[7:]

		id, err := common.ParseJWT(tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "权限不足"})
			c.Abort()
			return
		}

		admin := sdb.FindAdminByEvery("username", id)

		if admin == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "权限不足"})
			c.Abort()
			return
		}
		c.Set("user", admin)
		c.Next()

	}
}
