package login_log_panel

import (
	sdb "GinProject12/databases"
	"GinProject12/model"
	qqwey "GinProject12/util/resource"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func SaveLoginLog(c *gin.Context, err error) {
	ip := c.ClientIP()
	qwry, err := qqwey.NewQQwry()

	if err != nil {
		log.Println(err.Error())
		return
	}

	res := qwry.Find(ip)

	info := model.LoginLog{
		IsLogin:   err == nil,
		Ip:        res.IP,
		Area:      res.Area,
		LoginTime: time.Now().String(),
	}
	// 存入数据库
	if err = sdb.InsertLoginLog(info); err != nil {
		log.Println(err.Error())
		return
	}
}
