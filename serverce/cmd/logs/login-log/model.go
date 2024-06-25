package login_log

import "time"

type ResultQQwry struct {
	IP   string `json:"ip"`
	Area string `json:"area"`
}

type QQwry struct {
	Data   []byte
	Offset int64
}
type DeviceService struct{}

type sshFileItem struct {
	Name string
	Year int
}

type SSHService struct{}

type PageInfo struct {
	Page     int `json:"page" validate:"required,number"`
	PageSize int `json:"pageSize" validate:"required,number"`
}

type SearchSSHLog struct {
	PageInfo
	Info   string `json:"info"`
	Status string `json:"Status" validate:"required,oneof=Success Failed All"`
}
type SSHLog struct {
	Logs            []SSHHistory `json:"logs"`
	TotalCount      int          `json:"totalCount"`
	SuccessfulCount int          `json:"successfulCount"`
	FailedCount     int          `json:"failedCount"`
}

type SSHHistory struct {
	Date     time.Time `json:"date"`
	DateStr  string    `json:"dateStr"`
	Area     string    `json:"area"`
	User     string    `json:"user"`
	AuthMode string    `json:"authMode"`
	Address  string    `json:"address"`
	Port     string    `json:"port"`
	Status   string    `json:"status"`
	Message  string    `json:"message"`
}
