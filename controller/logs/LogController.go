package logs

import (
	"GinProject12/response"
	login_log "GinProject12/serverce/cmd/logs/login-log"
	"GinProject12/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// LoginLog 服务器登陆日志
// @Summary      服务器登陆日志
// @Description  服务器登陆日志
// @Tags         logs
// @Accept       json
// @Produce      json
// @Param        loginLog body login_log.SearchSSHLog ture "request"
// @Success      200  {object}  login_log.SSHLog
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /logs/login [POST]
func LoginLog(c *gin.Context) {

	var req login_log.SearchSSHLog

	if err := util.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	service := login_log.SSHService{}
	log, err := service.LoadLog(req)

	if err != nil {
		fmt.Println(err)
		response.Response(c, http.StatusInternalServerError, 500, nil, "服务器内部问题")
		return
	}

	response.Success(c, gin.H{"data": log}, "成功!")

}
