package monitor

import (
	sdb "GinProject12/databases"
	"GinProject12/model"
	condition "GinProject12/serverce/cmd/system/monitor"
	sysStatic "GinProject12/serverce/cmd/system/static"
	"GinProject12/util"
	"fmt"
	"log"
	"strings"
	"time"
)

/*

	数据库:
		创建时间    监控的硬件			监控的值							通知邮箱可以多个
				   cpu, 内存 ...		4大基本都用% 流量用上行下行 内存 读取 写入
*/

var (
	currentTime = time.Now()
)

func GetNowTime() string {
	currentTime = time.Now()
	return currentTime.Format("2006-01-02 15:04:05")
}

var AllMonitor = make(map[int]chan bool)

func (f *MonitorService) MonitorSwitch(req model.Monitor) error {
	// 插入数据库
	id, err := sdb.InsertMonitorInfo(req)
	if err != nil {
		return err
	}
	req.Id = id

	switch req.HardWare {
	case CPU:
		go MCpu(req)
	case MEMORY:
		go MMemory(req)
	case LOAD:
		go MLoad(req)
	case DISK:
		go MDisk(req)
	case NETWORK:
		go MNetwork(req)
	case IO:
		go MIo(req)
	}
	return nil
}

func (f *MonitorService) DelMonitor(id int) error {
	// 删除数据库
	if err := sdb.DelMonitorInfo(id); err != nil {
		return err
	}
	close(AllMonitor[id])
	return nil
}

func (f *MonitorService) SelectMonitor(page model.PageInfo) ([]model.Monitor, error) {
	return sdb.SelectMonitorPage(page)
}

func MCpu(req model.Monitor) {
	var cpu sysStatic.Cpu
	quit := make(chan bool)
	AllMonitor[req.Id] = quit

	for {
		select {
		case <-quit:
			return
		default:
			if err := cpu.GetCpuInfo(); err != nil {
				log.Println(err)
			}

			if req.Threshold <= cpu.CpuUsedPercent {
				emails := strings.Split(req.NotifyEmail, ",")
				for _, email := range emails {
					if err := util.SendEmailWarning(email, fmt.Sprintf("当前时间%s\ncpu阈值已达到%.2f您设定的是%.2f", GetNowTime(), cpu.CpuUsedPercent, req.Threshold)); err != nil {
						log.Println("发送失败")
						continue
					}
					time.Sleep(50 * time.Millisecond)
				}
				// 发送成功等待10分钟再检测
				time.Sleep(5 * time.Minute)
			}
			time.Sleep(3 * time.Second)

		}

	}

}

func MMemory(req model.Monitor) {
	var mem sysStatic.Memory
	// 数据库读取
	quit := make(chan bool)
	AllMonitor[req.Id] = quit
	for {
		select {
		case <-quit:
			return
		default:
			if err := mem.GetMemoryInfo(); err != nil {
				log.Println(err)
			}
			if mem.MUsedPercent >= req.Threshold {
				emails := strings.Split(req.NotifyEmail, ",")
				for _, email := range emails {
					if err := util.SendEmailWarning(email, fmt.Sprintf("当前时间%s\n内存阈值已达到%.2f您设定的是%.2f", GetNowTime(), mem.MUsedPercent, req.Threshold)); err != nil {
						log.Println("发送失败")
						continue
					}
					time.Sleep(50 * time.Millisecond)
				}
				// 发送成功等待10分钟再检测
				time.Sleep(5 * time.Minute)
			}
			time.Sleep(3 * time.Second)

		}

	}
}

func MLoad(req model.Monitor) {
	load := sysStatic.Load{}
	cpu := sysStatic.Cpu{}
	quit := make(chan bool)
	AllMonitor[req.Id] = quit
	for {
		select {
		case <-quit:
			return
		default:
			if err := load.GetLoadInfo(); err != nil {
				log.Println(err)
			}
			if err := cpu.GetCpuInfo(); err != nil {
				log.Println(err)
			}
			load.LoadUsagePercent = load.Load1 / (float64(cpu.NumCPU*2) * 0.75) * 100

			if load.LoadUsagePercent >= req.Threshold {
				emails := strings.Split(req.NotifyEmail, ",")
				for _, email := range emails {
					if err := util.SendEmailWarning(email, fmt.Sprintf("当前时间%s\n系统负载阈值已达到%.2f您设定的是%.2f", GetNowTime(), load.LoadUsagePercent, req.Threshold)); err != nil {
						log.Println("发送失败")
						continue
					}
					time.Sleep(50 * time.Millisecond)
				}
				// 发送成功等待10分钟再检测
				time.Sleep(5 * time.Minute)
			}
			time.Sleep(3 * time.Second)

		}

	}
}

func MDisk(req model.Monitor) {

	quit := make(chan bool)
	AllMonitor[req.Id] = quit

	var disk sysStatic.Disk

	for {
		select {
		case <-quit:
			return
		default:
			if err := disk.GetDiskInfo(); err != nil {
				log.Println(err)
			}
			if disk.UsedPercent >= req.Threshold {
				emails := strings.Split(req.NotifyEmail, ",")
				for _, email := range emails {
					if err := util.SendEmailWarning(email, fmt.Sprintf("当前时间%s\n磁盘阈值已达到%.2f您设定的是%.2f", GetNowTime(), disk.UsedPercent, req.Threshold)); err != nil {
						log.Println("发送失败")
						continue
					}
					time.Sleep(50 * time.Millisecond)
				}
				// 发送成功等待10分钟再检测
				time.Sleep(5 * time.Minute)
			}
			time.Sleep(3 * time.Second)

		}
	}
}

func MNetwork(req model.Monitor) {
	var net condition.Flow
	var t condition.Flow
	var flag = false
	quit := make(chan bool)
	AllMonitor[req.Id] = quit

	for {
		if req.Detail == "all" {
			if err := net.GetFlowAll(); err != nil {
				log.Println(err)
			}
		} else {
			if err := net.GetFlowInfoByName(req.Detail); err != nil {
				log.Println(err)
			}
		}
		DownRes := float64(net.BytesRecv-t.BytesRecv) / 1024 / 3
		UpRes := float64(net.BytesSent-t.BytesSent) / 1024 / 3

		if flag && DownRes >= req.Down {
			emails := strings.Split(req.NotifyEmail, ",")
			for _, email := range emails {
				if err := util.SendEmailWarning(email, fmt.Sprintf("当前时间%s\n网络下行阈值已达到%.2f 您设置的是%.2f", GetNowTime(), DownRes, req.Down)); err != nil {
					log.Println("发送失败")
					continue
				}
				time.Sleep(50 * time.Millisecond)
			}
			// 发送成功等待10分钟再检测
			time.Sleep(5 * time.Minute)
		}

		if flag && UpRes >= req.Up {
			emails := strings.Split(req.NotifyEmail, ",")
			for _, email := range emails {
				if err := util.SendEmailWarning(email, fmt.Sprintf("当前时间%s\n网络上行阈值已达到%.2f 您设置的是%.2f", GetNowTime(), UpRes, req.Up)); err != nil {
					log.Println("发送失败")
					continue
				}
				time.Sleep(50 * time.Millisecond)
			}
			// 发送成功等待10分钟再检测
			time.Sleep(5 * time.Minute)
		}

		flag = true
		t = net
		time.Sleep(3 * time.Second)
	}

}

func MIo(req model.Monitor) {
	var io condition.DiskIo
	var t condition.DiskIo
	var flag = false
	quit := make(chan bool)
	AllMonitor[req.Id] = quit

	for {
		if req.Detail == "all" {
			if err := io.GetDiskIoAll(); err != nil {
				log.Println(err)
			}
		} else {
			if err := io.GetDiskIoInfoByName(req.Detail); err != nil {
				log.Println(err)
			}
		}
		DownRes := float64(io.IoWriteBytes-t.IoWriteBytes) / 1024 / 1024 / 3
		UpRes := float64(io.IoReadBytes-t.IoReadBytes) / 1024 / 1024 / 3

		if flag && DownRes >= req.Down {
			emails := strings.Split(req.NotifyEmail, ",")
			for _, email := range emails {
				if err := util.SendEmailWarning(email, fmt.Sprintf("当前时间%s\n网络下行阈值已达到%.2f 您设置的是%.2f", GetNowTime(), DownRes, req.Down)); err != nil {
					log.Println("发送失败")
					continue
				}
				time.Sleep(50 * time.Millisecond)
			}
			// 发送成功等待10分钟再检测
			time.Sleep(5 * time.Minute)
		}

		if flag && UpRes >= req.Up {
			emails := strings.Split(req.NotifyEmail, ",")
			for _, email := range emails {
				if err := util.SendEmailWarning(email, fmt.Sprintf("当前时间%s\n网络上行阈值已达到%.2f 您设置的是%.2f", GetNowTime(), UpRes, req.Up)); err != nil {
					log.Println("发送失败")
					continue
				}
				time.Sleep(50 * time.Millisecond)
			}
			// 发送成功等待10分钟再检测
			time.Sleep(5 * time.Minute)
		}

		flag = true
		t = io
		time.Sleep(3 * time.Second)
	}

}
