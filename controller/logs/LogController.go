package logs

import (
	sdb "GinProject12/databases"
	"GinProject12/model"
	"GinProject12/response"
	login_log "GinProject12/serverce/cmd/logs/login-log-server"
	"GinProject12/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

// LoginLog 服务器登陆日志
// @Summary      服务器登陆日志
// @Description  服务器登陆日志
// @Tags         logs
// @Accept       json
// @Produce      json
// @Param        loginLog body login_log.SearchSSHLog ture "status Success 或 Failed 或 不写 不写代表全部"
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
		response.Response(c, http.StatusInternalServerError, 500, nil, "服务器内部问题")
		return
	}

	response.Success(c, gin.H{"data": log}, "成功!")

}

// PanelLogin 面板登陆日志
// @Summary      面板登陆日志
// @Description  面板登陆日志
// @Tags         logs
// @Accept       json
// @Produce      json
// @Param        page body model.PageInfo ture "request"
// @Success      200  {object} model.LoginLog
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /logs/panel [POST]
func PanelLogin(c *gin.Context) {
	var req model.PageInfo
	if err := util.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	page, err := sdb.LoginLogPage(req)

	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "服务器内部问题")
		return
	}

	response.Success(c, gin.H{"data": page}, "成功!")

}
