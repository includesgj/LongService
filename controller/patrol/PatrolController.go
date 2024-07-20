package patroltro

import (
	sdb "GinProject12/databases"
	"GinProject12/model"
	"GinProject12/response"
	"GinProject12/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// AddPatrol 添加任务信息
// @Summary      添加任务信息
// @Description 添加任务信息
// @Tags         patrol
// @Accept       json
// @Produce      json
// @Param        req body model.Patrol ture "需要预留几个空让普通用户填的具体看model.PatrolUser时间不需要传"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /patrol/add [POST]
func AddPatrol(c *gin.Context) {
	var req model.Patrol
	if err := util.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	req.CreateTime = time.Now()
	_, err := sdb.InsertPatrol(req)

	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}

	response.Success(c, nil, "成功!")
}

// DelPatrol 按照id删除任务
// @Summary      按照id删除任务
// @Description  按照id删除任务
// @Tags         patrol
// @Accept       json
// @Produce      json
// @Param        id query string false "后端传的id"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /patrol/del [GET]
func DelPatrol(c *gin.Context) {
	Sid := c.Query("id")
	id, err := strconv.Atoi(Sid)
	if err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, err.Error())
		return
	}
	if err = sdb.DelPatrol(id); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}
	response.Success(c, nil, "成功!")

}

// PatrolPage 查询任务信息
// @Summary      查询任务信息
// @Description 查询任务信息
// @Tags         patrol
// @Accept       json
// @Produce      json
// @Param        page body model.PageInfo ture "分页"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /patrol/sel [POST]
func PatrolPage(c *gin.Context) {
	var req model.PageInfo
	if err := util.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	list, err := sdb.SelectPatrol(req)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}

	response.Success(c, gin.H{"data": list}, "成功!")
}

// AddPatrolInfo 添加完成任务信息
// @Summary      添加完成任务信息
// @Description 添加完成任务信息
// @Tags         patrol
// @Accept       json
// @Produce      json
// @Param        req body model.PatrolUser ture "时间不需要传"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /patrol/user/add [POST]
func AddPatrolInfo(c *gin.Context) {
	var req model.PatrolUser
	if err := util.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	req.PatrolTime = time.Now()

	_, err := sdb.InsertPatrolUser(req)

	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}

	response.Success(c, nil, "成功!")
}

// DelPatrolInfo 按照id删除完成的任务
// @Summary      按照id删除完成的任务
// @Description  按照id删除完成的任务
// @Tags         patrol
// @Accept       json
// @Produce      json
// @Param        id query string false "后端传的id"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /monitor/user/del [GET]
func DelPatrolInfo(c *gin.Context) {
	Sid := c.Query("id")
	id, err := strconv.Atoi(Sid)
	if err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, err.Error())
		return
	}
	if err = sdb.DelPatrolUser(id); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}
	response.Success(c, nil, "成功!")

}

// PatrolInfoPage 查询完成的任务信息
// @Summary      查询完成的任务信息
// @Description 查询完成的任务信息
// @Tags         patrol
// @Accept       json
// @Produce      json
// @Param        page body model.PageInfo ture "分页"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /patrol/user/sel [POST]
func PatrolInfoPage(c *gin.Context) {
	var req model.PageInfo
	if err := util.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	list, err := sdb.SelectPatrolUser(req)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}

	response.Success(c, gin.H{"data": list}, "成功!")
}
