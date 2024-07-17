package condition

import (
	"GinProject12/util"
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/net"
)

// GetFlowAll 获取网络的全部 all
func (p *Flow) GetFlowAll() error {
	all, err := net.IOCounters(false)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	p.Sent, p.SentUnit = util.MaxUnit(all[0].BytesSent)
	p.Recv, p.RecvUnit = util.MaxUnit(all[0].BytesRecv)
	p.BytesRecv, p.BytesSent = all[0].BytesRecv, all[0].BytesSent

	return nil
}

// GetFlowInfoByName 获取全部的详细信息
func (p *Flow) GetFlowInfoByName(name string) error {

	all, err := net.IOCounters(true)

	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	for _, item := range all {
		if item.Name == name {
			p.Sent, p.SentUnit = util.MaxUnit(item.BytesSent)
			p.Recv, p.RecvUnit = util.MaxUnit(item.BytesRecv)
			p.BytesRecv, p.BytesSent = item.BytesRecv, item.BytesSent
		}

	}

	return nil
}

func (f *NetOrDiskName) GetNetOrDiskName(name string) error {
	if name == "net" {
		all, err := net.IOCounters(true)

		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
		for _, item := range all {
			f.Name = append(f.Name, item.Name)
		}

	} else if name == "io" {
		counters, err := disk.IOCounters()
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		for _, item := range counters {
			f.Name = append(f.Name, item.Name)
		}
	} else {
		return errors.New("没有该参数")
	}
	return nil
}

// GetDiskIoInfoByName 按名字获取磁盘io的详细信息
func (p *DiskIo) GetDiskIoInfoByName(name string) error {
	counters, err := disk.IOCounters()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	var flag = false

	for _, o := range counters {

		if o.Name == name {
			flag = true
			p.IoWriteBytes = o.WriteBytes
			p.IoReadBytes = o.ReadBytes
			p.IoWriterTime = o.WriteTime
			p.IoReadTime = o.ReadTime
			p.IoCount = o.WriteCount + o.ReadCount
		}

	}

	if !flag {
		return errors.New("参数错误")
	}

	return nil

}

// GetDiskIoAll 获取磁盘io的All信息
func (p *DiskIo) GetDiskIoAll() error {
	counters, err := disk.IOCounters()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	totalReadBytes := uint64(0)
	totalWriteBytes := uint64(0)
	totalCount := uint64(0)
	totalWriteTime := uint64(0)
	totalReadTime := uint64(0)

	for _, ioCounter := range counters {
		totalReadBytes += ioCounter.ReadBytes
		totalWriteBytes += ioCounter.WriteBytes
		totalCount += ioCounter.ReadCount + ioCounter.WriteCount
		totalWriteTime += ioCounter.WriteTime
		totalReadTime += ioCounter.ReadTime
	}

	p.IoWriteBytes = totalWriteBytes
	p.IoReadBytes = totalReadBytes
	p.IoWriterTime = totalWriteTime
	p.IoReadTime = totalReadTime
	p.IoCount = totalCount

	return nil
}
