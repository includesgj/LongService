package monitortro

import (
	"GinProject12/model"
	"GinProject12/response"
	"GinProject12/serverce/cmd/monitor"
	"GinProject12/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var monitorService monitor.MonitorService

// Monitor 添加监控系统信息
// @Summary      添加监控系统信息
// @Description 添加监控系统信息
// @Tags         monitor
// @Accept       json
// @Produce      json
// @Param        req body model.Monitor ture "需要用户 hardware是检测的硬件 detail是如果是io或net才需要传要不是all要不是单个名称 Threshold是百分比用户设定的 up是io或net的上行或上传 down反过来 notifyEmail是要通知邮箱可以多个传过来要用,分割"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /monitor/add [POST]
func Monitor(c *gin.Context) {
	var req model.Monitor
	if err := util.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	if err := monitorService.MonitorAdd(req); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}

	response.Success(c, nil, "成功!")

}

// DelMonitor 按照id删除监控
// @Summary      按照id删除监控
// @Description  按照id删除监控
// @Tags         monitor
// @Accept       json
// @Produce      json
// @Param        id query string false "后端传的id"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /monitor/del [GET]
func DelMonitor(c *gin.Context) {
	Sid := c.Query("id")
	id, err := strconv.Atoi(Sid)
	if err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, err.Error())
		return
	}

	if err = monitorService.DelMonitor(id); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}

	response.Success(c, nil, "成功!")
}

// Monitor 查询监控系统信息
// @Summary      查询监控系统信息
// @Description 查询监控系统信息
// @Tags         monitor
// @Accept       json
// @Produce      json
// @Param        page body model.PageInfo ture "分页"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /monitor/sel [POST]
func SelectMonitor(c *gin.Context) {
	var req model.PageInfo
	if err := util.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	list, err := monitorService.SelectMonitor(req)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}

	response.Success(c, gin.H{"data": list}, "成功!")

}
