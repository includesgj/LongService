package systemtro

import (
	"GinProject12/response"
	"GinProject12/serverce/cmd/system/mes"
	"GinProject12/serverce/cmd/system/monitor"
	sysStatic "GinProject12/serverce/cmd/system/static"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SysInfo 服务器硬件信息
// @Summary      服务器硬件信息
// @Description  服务器硬件信息
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object}  sysMes.LocalInfo
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /sys/info [GET]
func SysInfo(c *gin.Context) {
	var SInfo = &sysMes.LocalInfo{}
	if err := SInfo.GetSysInfo(); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "获取失败")
	}

	response.Success(c, gin.H{"data": SInfo}, "成功")

}

// SysStatic 获取服务器状态
// @Summary     获取服务器状态
// @Description  获取服务器状态
// @Tags         system
// @Accept       json
// @Produce      json
// @Success      200  {object} sysStatic.SysStatus
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /sys/static [GET]
func SysStatic(c *gin.Context) {

	cpu := sysStatic.Cpu{}
	memory := sysStatic.Memory{}
	load := sysStatic.Load{}
	disk := sysStatic.Disk{}

	if err := cpu.GetCpuInfo(); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "获取失败")
	}
	if err := memory.GetMemoryInfo(); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "获取失败")
	}
	if err := load.GetLoadInfo(); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "获取失败")
	}
	if err := disk.GetDiskInfo(); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "获取失败")
	}

	// 这里需要cpu个数来计算总负载率
	load.LoadUsagePercent = load.Load1 / (float64(cpu.NumCPU*2) * 0.75) * 100
	status := sysStatic.SysStatus{Cpu: cpu, Memory: memory, Load: load, Disk: disk}

	response.Success(c, gin.H{"data": status}, "成功")

}

// SysMonitorNet 获取网络监控信息 name表示获取具体哪个
// @Summary      获取网络监控信息
// @Description  获取网络监控信息
// @Tags         system
// @Accept       json
// @Produce      json
// @Param        name query string false "该值表示要监控的网络名称All表示全部"
// @Success      200 {object} condition.Flow
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /sys/net [GET]
func SysMonitorNet(c *gin.Context) {
	name := c.Query("name")
	flow := &condition.Flow{}

	if name != "all" {
		if err := flow.GetFlowInfoByName(name); err != nil {
			response.Response(c, http.StatusInternalServerError, 500, nil, "获取失败")
		}
	} else {
		if err := flow.GetFlowAll(); err != nil {
			response.Response(c, http.StatusInternalServerError, 500, nil, "获取失败")
		}
	}
	response.Success(c, gin.H{"data": flow}, "成功")

}

// SysMonitorIo 获取io监控信息
// @Summary      获取io监控信息
// @Description  获取io监控信息
// @Tags         system
// @Accept       json
// @Produce      json
// @Param        name query string false "该值表示要监控的磁盘名称All表示全部"
// @Success      200  {object}  condition.DiskIo
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /sys/io [GET]
func SysMonitorIo(c *gin.Context) {
	name := c.Query("name")
	io := &condition.DiskIo{}

	if name != "" {
		if err := io.GetDiskIoInfoByName(name); err != nil {
			response.Response(c, http.StatusBadRequest, 400, nil, "获取失败")
		}
	} else {
		if err := io.GetDiskIoAll(); err != nil {
			response.Response(c, http.StatusInternalServerError, 500, nil, "获取失败")
		}
	}
	response.Success(c, gin.H{"data": io}, "成功")
}

// GetNetOrDiskName 获取io或net名称
// @Summary      获取io或net名称
// @Description  获取io或net名称
// @Tags         system
// @Accept       json
// @Produce      json
// @Param        name query string false "io 或者是 net"
// @Success      200  {object}  condition.NetOrDiskName
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /sys/name [GET]
func GetNetOrDiskName(c *gin.Context) {
	name := c.Query("name")

	var s condition.NetOrDiskName
	if err := s.GetNetOrDiskName(name); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "获取失败")
	}
	response.Success(c, gin.H{"data": s}, "成功")

}
