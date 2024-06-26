package login_log

import (
	qqwey "GinProject12/serverce/cmd/logs/login-log/resource"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	StatusSuccess = "Success"
	StatusFailed  = "Failed"
)

const (
	indexLen      = 7
	redirectMode1 = 0x01
	redirectMode2 = 0x02
)

var (
	ErrCmdTimeout = "ErrCmdTimeout"
)

func handleGunzip(path string) error {
	if _, err := Execf("gunzip %s", path); err != nil {
		return err
	}
	return nil
}

func Execf(cmdStr string, a ...interface{}) (string, error) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf(cmdStr, a...))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return stdout.String(), nil
}

func sortFileList(fileNames []sshFileItem) []sshFileItem {
	if len(fileNames) < 2 {
		return fileNames
	}
	if strings.HasPrefix(path.Base(fileNames[0].Name), "secure") {
		var itemFile []sshFileItem
		sort.Slice(fileNames, func(i, j int) bool {
			return fileNames[i].Name > fileNames[j].Name
		})
		itemFile = append(itemFile, fileNames[len(fileNames)-1])
		itemFile = append(itemFile, fileNames[:len(fileNames)-1]...)
		return itemFile
	}
	sort.Slice(fileNames, func(i, j int) bool {
		return fileNames[i].Name < fileNames[j].Name
	})
	return fileNames
}

func Exec(cmdStr string) (string, error) {
	return ExecWithTimeOut(cmdStr, 20*time.Second)
}

func (u *DeviceService) LoadTimeZone() ([]string, error) {
	std, err := Exec("timedatectl list-timezones")
	if err != nil {
		return []string{}, err
	}
	return strings.Split(std, "\n"), nil
}

func ExecWithTimeOut(cmdStr string, timeout time.Duration) (string, error) {
	cmd := exec.Command("bash", "-c", cmdStr)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return "", err
	}
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	after := time.After(timeout)
	select {
	case <-after:
		_ = cmd.Process.Kill()
		return "", errors.New(ErrCmdTimeout)
	case err := <-done:
		if err != nil {
			return "", err
		}
	}

	return stdout.String(), nil
}

func LoadTimeZone() string {
	loc := time.Now().Location()
	if _, err := time.LoadLocation(loc.String()); err != nil {
		return "Asia/Shanghai"
	}
	return loc.String()
}

func NewQQwry() (*QQwry, error) {
	IpCommonDictionary := qqwey.QQwryByte
	return &QQwry{Data: IpCommonDictionary}, nil
}

func analyzeDateStr(parts []string) (int, string) {
	t, err := time.Parse("2006-01-02T15:04:05.999999-07:00", parts[0])
	if err != nil {
		if len(parts) < 14 {
			return 0, ""
		}
		return 2, fmt.Sprintf("%s %s %s", parts[0], parts[1], parts[2])
	}
	if len(parts) < 12 {
		return 0, ""
	}
	return 0, t.Format("2006 Jan 2 15:04:05")
}

func loadFailedSecureDatas(line string) SSHHistory {
	var data SSHHistory
	// 空格分割
	parts := strings.Fields(line)
	index, dataStr := analyzeDateStr(parts)
	if dataStr == "" {
		return data
	}
	data.DateStr = dataStr
	if strings.Contains(line, " invalid ") {
		data.AuthMode = parts[4+index]
		index += 2
	} else {
		data.AuthMode = parts[4+index]
	}
	data.User = parts[6+index]
	data.Address = parts[8+index]
	data.Port = parts[10+index]
	data.Status = StatusFailed
	if strings.Contains(line, ": ") {
		data.Message = strings.Split(line, ": ")[1]
	}
	return data
}

func loadDate(currentYear int, DateStr string, nyc *time.Location) time.Time {
	itemDate, err := time.ParseInLocation("2006 Jan 2 15:04:05", fmt.Sprintf("%d %s", currentYear, DateStr), nyc)
	if err != nil {
		itemDate, _ = time.ParseInLocation("2006 Jan 2 15:04:05", DateStr, nyc)
	}
	return itemDate
}

func loadFailedAuthDatas(line string) SSHHistory {
	var data SSHHistory
	parts := strings.Fields(line)
	index, dataStr := analyzeDateStr(parts)
	if dataStr == "" {
		return data
	}
	data.DateStr = dataStr
	if index == 2 {
		data.User = parts[10]
	} else {
		data.User = parts[7]
	}
	data.AuthMode = parts[6+index]
	data.Address = parts[9+index]
	data.Port = parts[11+index]
	data.Status = StatusFailed
	if strings.Contains(line, ": ") {
		data.Message = strings.Split(line, ": ")[1]
	}
	return data
}

func loadSuccessDatas(line string) SSHHistory {
	var data SSHHistory
	parts := strings.Fields(line)
	index, dataStr := analyzeDateStr(parts)
	if dataStr == "" {
		return data
	}
	data.DateStr = dataStr
	data.AuthMode = parts[4+index]
	data.User = parts[6+index]
	data.Address = parts[8+index]
	data.Port = parts[10+index]
	data.Status = StatusSuccess
	return data
}

func loadSSHData(command string, showCountFrom, showCountTo, currentYear int, qqWry *QQwry, nyc *time.Location) ([]SSHHistory, int, int) {
	var (
		datas        []SSHHistory
		successCount int
		failedCount  int
	)
	stdout2, err := Exec(command)
	if err != nil {
		return datas, 0, 0
	}
	lines := strings.Split(string(stdout2), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		var itemData SSHHistory
		switch {
		case strings.Contains(lines[i], "Failed password for"):
			itemData = loadFailedSecureDatas(lines[i])
			if len(itemData.Address) != 0 {
				if successCount+failedCount >= showCountFrom && successCount+failedCount < showCountTo {
					itemData.Area = qqWry.Find(itemData.Address).Area
					itemData.Date = loadDate(currentYear, itemData.DateStr, nyc)
					datas = append(datas, itemData)
				}
				failedCount++
			}
		case strings.Contains(lines[i], "Connection closed by authenticating user"):
			itemData = loadFailedAuthDatas(lines[i])
			if len(itemData.Address) != 0 {
				if successCount+failedCount >= showCountFrom && successCount+failedCount < showCountTo {
					itemData.Area = qqWry.Find(itemData.Address).Area
					itemData.Date = loadDate(currentYear, itemData.DateStr, nyc)
					datas = append(datas, itemData)
				}
				failedCount++
			}
		case strings.Contains(lines[i], "Accepted "):
			itemData = loadSuccessDatas(lines[i])
			if len(itemData.Address) != 0 {
				if successCount+failedCount >= showCountFrom && successCount+failedCount < showCountTo {
					itemData.Area = qqWry.Find(itemData.Address).Area
					itemData.Date = loadDate(currentYear, itemData.DateStr, nyc)
					datas = append(datas, itemData)
				}
				successCount++
			}
		}
	}
	return datas, successCount, failedCount
}

func (u *SSHService) LoadLog(req SearchSSHLog) (*SSHLog, error) {
	var fileList []sshFileItem
	var data SSHLog
	baseDir := "/var/log"
	// 查询 /var/log 目录下secure 或者 auth 开头的文件如果是.gz就解压最后都放到fileList上
	if err := filepath.Walk(baseDir, func(pathItem string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasPrefix(info.Name(), "secure") || strings.HasPrefix(info.Name(), "auth")) {
			if !strings.HasSuffix(info.Name(), ".gz") {
				fileList = append(fileList, sshFileItem{Name: pathItem, Year: info.ModTime().Year()})
				return nil
			}
			itemFileName := strings.TrimSuffix(pathItem, ".gz")
			if _, err := os.Stat(itemFileName); err != nil && os.IsNotExist(err) {
				if err := handleGunzip(pathItem); err == nil {
					fileList = append(fileList, sshFileItem{Name: itemFileName, Year: info.ModTime().Year()})
				}
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	// 排序
	fileList = sortFileList(fileList)

	// info 是前端传过来的 正则用的
	command := ""
	if len(req.Info) != 0 {
		command = fmt.Sprintf(" | grep '%s'", req.Info)
	}

	showCountFrom := (req.Page - 1) * req.PageSize
	showCountTo := req.Page * req.PageSize
	// 加载当地时区
	nyc, _ := time.LoadLocation(LoadTimeZone())
	qqWry, err := NewQQwry()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for _, file := range fileList {
		commandItem := ""
		if strings.HasPrefix(path.Base(file.Name), "secure") {
			switch req.Status {
			case StatusSuccess:
				commandItem = fmt.Sprintf("cat %s | grep -a Accepted %s", file.Name, command)
			case StatusFailed:
				commandItem = fmt.Sprintf("cat %s | grep -a 'Failed password for' %s", file.Name, command)
			default:
				commandItem = fmt.Sprintf("cat %s | grep -aE '(Failed password for|Accepted)' %s", file.Name, command)
			}
		}
		if strings.HasPrefix(path.Base(file.Name), "auth.log") {
			switch req.Status {
			case StatusSuccess:
				commandItem = fmt.Sprintf("cat %s | grep -a Accepted %s", file.Name, command)
			case StatusFailed:
				commandItem = fmt.Sprintf("cat %s | grep -aE 'Failed password for|Connection closed by authenticating user' %s", file.Name, command)
			default:
				commandItem = fmt.Sprintf("cat %s | grep -aE \"(Failed password for|Connection closed by authenticating user|Accepted)\" %s", file.Name, command)
			}
		}
		dataItem, successCount, failedCount := loadSSHData(commandItem, showCountFrom, showCountTo, file.Year, qqWry, nyc)
		data.FailedCount += failedCount
		data.TotalCount += successCount + failedCount
		showCountFrom = showCountFrom - (successCount + failedCount)
		showCountTo = showCountTo - (successCount + failedCount)
		data.Logs = append(data.Logs, dataItem...)
	}

	data.SuccessfulCount = data.TotalCount - data.FailedCount
	return &data, nil
}

// readData 从文件中读取数据
func (q *QQwry) readData(num int, offset ...int64) (rs []byte) {
	if len(offset) > 0 {
		q.setOffset(offset[0])
	}
	nums := int64(num)
	end := q.Offset + nums
	dataNum := int64(len(q.Data))
	if q.Offset > dataNum {
		return nil
	}

	if end > dataNum {
		end = dataNum
	}
	rs = q.Data[q.Offset:end]
	q.Offset = end
	return
}

// setOffset 设置偏移量
func (q *QQwry) setOffset(offset int64) {
	q.Offset = offset
}

// Find ip地址查询对应归属地信息
func (q *QQwry) Find(ip string) (res ResultQQwry) {
	res = ResultQQwry{}
	res.IP = ip
	if strings.Count(ip, ".") != 3 {
		return res
	}
	offset := q.searchIndex(binary.BigEndian.Uint32(net.ParseIP(ip).To4()))
	if offset <= 0 {
		return
	}

	var area []byte
	mode := q.readMode(offset + 4)
	if mode == redirectMode1 {
		countryOffset := q.readUInt24()
		mode = q.readMode(countryOffset)
		if mode == redirectMode2 {
			c := q.readUInt24()
			area = q.readString(c)
		} else {
			area = q.readString(countryOffset)
		}
	} else if mode == redirectMode2 {
		countryOffset := q.readUInt24()
		area = q.readString(countryOffset)
	} else {
		area = q.readString(offset + 4)
	}

	enc := simplifiedchinese.GBK.NewDecoder()
	res.Area, _ = enc.String(string(area))

	return
}

// readMode 获取偏移值类型
func (q *QQwry) readMode(offset uint32) byte {
	mode := q.readData(1, int64(offset))
	return mode[0]
}

// readString 获取字符串
func (q *QQwry) readString(offset uint32) []byte {
	q.setOffset(int64(offset))
	data := make([]byte, 0, 30)
	for {
		buf := q.readData(1)
		if buf[0] == 0 {
			break
		}
		data = append(data, buf[0])
	}
	return data
}

// searchIndex 查找索引位置
func (q *QQwry) searchIndex(ip uint32) uint32 {
	header := q.readData(8, 0)

	start := binary.LittleEndian.Uint32(header[:4])
	end := binary.LittleEndian.Uint32(header[4:])

	for {
		mid := q.getMiddleOffset(start, end)
		buf := q.readData(indexLen, int64(mid))
		_ip := binary.LittleEndian.Uint32(buf[:4])

		if end-start == indexLen {
			offset := byteToUInt32(buf[4:])
			buf = q.readData(indexLen)
			if ip < binary.LittleEndian.Uint32(buf[:4]) {
				return offset
			}
			return 0
		}

		if _ip > ip {
			end = mid
		} else if _ip < ip {
			start = mid
		} else if _ip == ip {
			return byteToUInt32(buf[4:])
		}
	}
}

// readUInt24
func (q *QQwry) readUInt24() uint32 {
	buf := q.readData(3)
	return byteToUInt32(buf)
}

// getMiddleOffset
func (q *QQwry) getMiddleOffset(start uint32, end uint32) uint32 {
	records := ((end - start) / indexLen) >> 1
	return start + records*indexLen
}

// byteToUInt32 将 byte 转换为uint32
func byteToUInt32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}