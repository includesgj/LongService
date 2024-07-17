package monitor

/*
MonitorReq
createtime 创建时间
hardware  监视的是什么
Detail 如果是网络或者是磁盘io就有具体的哪一个或者是all全部
threshold 百分比 如果是监视 cpu 内存 负载 磁盘使用率
up 网络的上行 或 磁盘io的读取
down 网络的下行 或 磁盘io的写入
notifyemail通知邮箱
*/

const (
	CPU     = "cpu"
	MEMORY  = "memory"
	LOAD    = "load"
	DISK    = "disk"
	NETWORK = "net"
	IO      = "io"
)

type MonitorService struct {
}
