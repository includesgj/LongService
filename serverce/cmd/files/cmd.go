package files

import (
	"errors"
	"fmt"
	"os/exec"
)

func Which(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func checkCmdAvailability(cmdStr string) error {
	if Which(cmdStr) {
		return nil
	}
	return errors.New("该命令为找到" + cmdStr)
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
