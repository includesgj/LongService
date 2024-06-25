package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Response(c *gin.Context, httpStatic int, code int, data gin.H, msg string) {
	c.JSON(httpStatic, gin.H{"code": code, "data": data, "mes": msg})
	c.Abort()
}

func Success(c *gin.Context, data gin.H, msg string) {
	Response(c, http.StatusOK, 200, data, msg)
}

func Fail(c *gin.Context, data gin.H, msg string) {
	Response(c, http.StatusBadRequest, 400, data, msg)
}
