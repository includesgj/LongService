package sysMes

/*
LocalInfo
发行版本: platform
主机名: hostname
内核版本: kernelVersion
系统类型: kernelArch
启动时间和运行时间 在 uptime上去算
*/
type LocalInfo struct {
	Platform  string `json:"platform"`
	HostName  string `json:"hostName"`
	KernelV   string `json:"kernelV"`
	SysType   string `json:"sysType"`
	StartTime string `json:"startTime"`
	RunTime   string `json:"runTime"`
}
