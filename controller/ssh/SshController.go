package sshtro

import (
	"GinProject12/response"
	sshs "GinProject12/serverce/cmd/ssh"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var server sshs.SSHService

// SshService 修改ssh服务
// @Summary      修改ssh服务
// @Description  修改ssh服务 operate有两种值 1, enable (开启) 2, disable (关闭)
// @Tags         email
// @Accept       json
// @Produce      json
// @Param        operate query string ture "operate"
// @Success      200
// @Failure      201
// @Failure      404
// @Failure      500
// @Router       /ssh/operate [GET]
func SshService(c *gin.Context) {
	operate := c.Query("operate")
	if operate != "disable" && operate != "enable" {
		response.Response(c, http.StatusBadRequest, 400, nil, "参数不正确")
		return
	}

	if err := server.OperationSsh(operate); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}
	response.Success(c, nil, fmt.Sprintf("%s成功", operate))
}
