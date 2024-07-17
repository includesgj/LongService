package sysStatic

import (
	"GinProject12/util"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"strconv"
)

func (p *Cpu) GetCpuInfo() error {

	cpuPercent, err := cpu.Percent(0, true)

	if err != nil {
		fmt.Println("获取CPU信息错误", err)
		return err
	}

	totalPercent, err := cpu.Percent(0, false)
	if err != nil {
		return err
	}

	if len(totalPercent) == 1 {
		p.CpuUsedPercent = totalPercent[0]
		p.CpuUsed = p.CpuUsedPercent * 0.01 * float64(p.NumCPU)
	}

	// cpu的使用率
	for _, c := range cpuPercent {
		p.Percentage = append(p.Percentage, float32(c))
	}

	// cpu型号
	cpuInfo, err := cpu.Info()
	if err != nil {
		fmt.Println("获取CPU信息错误", err)
		return err
	}
	if len(cpuInfo) > 0 {

		for _, e := range cpuInfo {
			cnt, err := strconv.Atoi(e.CoreID)
			if err != nil {
				fmt.Println("解析错误", err)
				return err
			}

			if uint(cnt) > p.NumCPU {
				p.NumCPU = uint(cnt)
			}

		}
		p.NumLogicCpu = uint(len(cpuInfo))
		p.NumCPU++
		p.Model = cpuInfo[0].ModelName
	}

	return nil

}

func bytesToMB(bytes uint64) float64 {
	return float64(bytes) / 1024 / 1024
}

func (p *Memory) GetMemoryInfo() error {
	memInfo, err := mem.VirtualMemory()

	if err != nil {
		fmt.Println("查看memory错误", err)
		return err
	}

	p.MTotal = bytesToMB(memInfo.Total)
	p.Available = bytesToMB(memInfo.Available)
	p.MUsed = bytesToMB(memInfo.Used)
	p.MFree = bytesToMB(memInfo.Free)
	p.MUsedPercent = memInfo.UsedPercent
	p.SwapTotal = bytesToMB(memInfo.SwapTotal)
	p.SwapFree = bytesToMB(memInfo.SwapFree)
	p.SwapPercent = p.SwapUsed / p.SwapTotal

	return nil
}

func (p *Load) GetLoadInfo() error {
	avg, err := load.Avg()
	if err != nil {
		fmt.Println("获取load错误", err)
		return err
	}

	p.Load1 = avg.Load1
	p.Load5 = avg.Load5
	p.Load15 = avg.Load15

	return nil
}

func (p *Disk) GetDiskInfo() error {

	partitions, err := disk.Partitions(false)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	p.Path = partitions[0].Mountpoint
	p.Device = partitions[0].Device
	p.FsType = partitions[0].Fstype

	fsInfo, err := disk.Usage(p.Path)
	if err != nil {
		fmt.Println("获取文件系统信息出错:", err)
		return err
	}

	p.UsedPercent = fsInfo.UsedPercent
	p.InodesTotal = fsInfo.InodesTotal
	p.InodesUsed = fsInfo.InodesUsed
	p.InodesFree = fsInfo.InodesFree
	p.InodesUsedPercent = fsInfo.InodesUsedPercent

	total, s := util.MaxUnit(fsInfo.Total)
	p.Total = fmt.Sprintf("%.2f %s", total, s)
	used, s := util.MaxUnit(fsInfo.Used)
	p.Used = fmt.Sprintf("%.2f %s", used, s)
	free, s := util.MaxUnit(fsInfo.Free)
	p.Free = fmt.Sprintf("%.2f %s", free, s)

	return nil
}
