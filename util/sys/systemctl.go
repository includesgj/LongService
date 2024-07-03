package sys

import (
	"os/exec"
	"strings"
)

func SystemCtlCmd(cmd ...string) (string, error) {
	res := exec.Command("systemctl", cmd...)
	output, err := res.CombinedOutput()
	if err != nil {
		return string(output), err
	}
	return string(output), err
}

func IsExist(c string) (bool, error) {
	res, err := SystemCtlCmd("is-enabled", c)
	if err != nil {
		// 如果没有启动也是会返回错误的
		if strings.Contains(res, "disabled") {
			return true, err
		}
		return false, err
	}
	// 此处就是启用了
	return true, nil
}

func RemoteMode() (string, error) {
	var err error
	var ok bool
	if ok, err = IsExist("sshd"); ok {
		return "sshd", nil
	}

	if ok, err = IsExist("ssh"); ok {
		return "ssh", nil
	}

	return "", err
}
