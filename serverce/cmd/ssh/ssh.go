package sshs

import (
	"GinProject12/util/cmd"
	"GinProject12/util/sys"
	"errors"
	"fmt"
	"strings"
)

type SSHService struct {
}

func (f *SSHService) OperationSsh(operate string) error {
	mode, err := sys.RemoteMode()
	if err != nil {
		return err
	}

	var name string

	if operate == "disable" || operate == "enable" {
		name = mode + ".service"
	}

	sudo := cmd.ChickSudoCmd()

	exec, err := cmd.Exec("%s systemctl %s %s", sudo, operate, name)

	if err != nil {
		// 操作的服务存在别名或者软链接
		if mode != "ssh" && strings.Contains(exec, "alias name or linked unit file") {
			exec, err = cmd.Exec("%s systemctl %s ssh", sudo, operate)
			if err != nil {
				return errors.New("alias name or linked unit file err: " + err.Error())
			}
		}
		return fmt.Errorf("%s -> %s 失败 err ->%v", name, operate, err)
	}

	return nil
}
