package monitor

/*
Flow
Type 代表是全部还是一个一个的详细
BytesSent BytesRecv byte为单位的
Unit 单位
Sent Recv 最大单位
*/
type Flow struct {
	BytesSent uint64  `json:"bytesSent"`
	BytesRecv uint64  `json:"bytesRecv"`
	SentUnit  string  `json:"sentUnit"`
	RecvUnit  string  `json:"recvUnit"`
	Sent      float64 `json:"sent"`
	Recv      float64 `json:"recv"`
}

type DiskIo struct {
	IoCount      uint64 `json:"ioCount"`
	IoWriteBytes uint64 `json:"ioWriteBytes"`
	IoReadBytes  uint64 `json:"ioReadBytes"`
	IoWriterTime uint64 `json:"ioWriterTime"`
	IoReadTime   uint64 `json:"ioReadTime"`
}
