package sysStatic

type SysStatus struct {
	Cpu
	Memory
	Disk
	Load
}

type Cpu struct {
	NumCPU         uint      `json:"numCPU"`
	NumLogicCpu    uint      `json:"numLogicCpu"`
	Model          string    `json:"model"`
	CpuUsedPercent float64   `json:"cpuUsedPercent"`
	CpuUsed        float64   `json:"cpuUsed"`
	Percentage     []float32 `json:"percentage"`
}

type Memory struct {
	MTotal       float64 `json:"mTotal"`
	Available    float64 `json:"available"`
	MUsed        float64 `json:"mUsed"`
	MFree        float64 `json:"mFree"`
	MUsedPercent float64 `json:"mUsedPercent"`
	SwapTotal    float64 `json:"swapTotal"`
	SwapUsed     float64 `json:"swapUsed"`
	SwapFree     float64 `json:"swapFree"`
	SwapPercent  float64 `json:"swapPercent"`
}

type Load struct {
	Load1            float64 `json:"load1"`
	Load5            float64 `json:"load5"`
	Load15           float64 `json:"load15"`
	LoadUsagePercent float64 `json:"loadUsagePercent"`
}

type Disk struct {
	Path              string  `json:"path"`
	FsType            string  `json:"fstype"`
	Device            string  `json:"device"`
	Total             string  `json:"total"`
	Free              string  `json:"free"`
	Used              string  `json:"used"`
	UsedPercent       float64 `json:"usedPercent"`
	InodesTotal       uint64  `json:"inodesTotal"`
	InodesUsed        uint64  `json:"inodesUsed"`
	InodesFree        uint64  `json:"inodesFree"`
	InodesUsedPercent float64 `json:"inodesUsedPercent"`
}
