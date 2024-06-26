package files

import (
	"fmt"
)

type TarArchiver struct {
	Cmd          string
	CompressType CompressType
}

type ShellArchiver interface {
	Extract(filePath, dstDir string) error
	Compress(sourcePaths []string, dstFile string) error
}

func NewTarArchiver(compressType CompressType) ShellArchiver {
	return &TarArchiver{
		Cmd:          "tar",
		CompressType: compressType,
	}
}

func (t TarArchiver) Extract(FilePath string, dstDir string) error {
	return ExecCmd(fmt.Sprintf("%s %s %s -C %s", t.Cmd, t.getOptionStr("extract"), FilePath, dstDir))
}

func (t TarArchiver) Compress(sourcePaths []string, dstFile string) error {
	return nil
}

func (t TarArchiver) getOptionStr(Option string) string {
	switch t.CompressType {
	case Tar:
		if Option == "compress" {
			return "cvf"
		} else {
			return "xf"
		}
	}
	return ""
}
