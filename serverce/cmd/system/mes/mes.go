package sysMes

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/host"
	"log"
	"time"
)

// GetSysInfo 系统信息
func (p *LocalInfo) GetSysInfo() error {
	hostInfo, err := host.Info()
	if err != nil {
		log.Printf("Failed to get host info: %v", err)
		return err
	}

	p.Platform = hostInfo.Platform
	p.HostName = hostInfo.Hostname
	p.KernelV = hostInfo.KernelVersion
	p.SysType = hostInfo.KernelArch

	startTimeUnix := int64(hostInfo.BootTime)
	startTime := time.Unix(startTimeUnix, 0)

	// 格式化时间
	layout := "2006-01-02 15:04:05" // 固定时间格式
	p.StartTime = startTime.Format(layout)

	seconds := hostInfo.Uptime % 60
	minutes := hostInfo.Uptime / 60 % 60
	hours := hostInfo.Uptime / 60 / 60 % 24
	days := hostInfo.Uptime / 60 / 60 / 24

	p.RunTime = fmt.Sprintf("%d天 %02d小时%02d分钟%02d秒", days, hours, minutes, seconds)

	return nil
}
