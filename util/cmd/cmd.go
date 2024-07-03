package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

func Which(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func CheckCmdAvailability(cmdStr string) error {
	if Which(cmdStr) {
		return nil
	}
	return errors.New("该命令未找到" + cmdStr)
}

func ExecCmdWithDir(cmdStr, workDir string) error {
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error : %v, output: %s", err, output)
	}
	return nil
}

func ExecCmd(cmdStr string) error {
	cmd := exec.Command("bash", "-c", cmdStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error : %v, output: %s", err, output)
	}
	return nil
}

func Exec(cmdStr string, a ...interface{}) (string, error) {
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

// 检测有没有sudo权限
func ChickSudoCmd() string {
	// w是当前登陆系统的用户信息
	res := exec.Command("sudo", "-u", "w")

	if err := res.Run(); err != nil {
		return ""
	}
	return "sudo "
}
