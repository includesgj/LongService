package futil

import (
	"fmt"
	"github.com/spf13/afero"
	"io/fs"
	"os"
	"path/filepath"
)

type FileOp struct {
	Fs afero.Fs
}

func NewFileOp() FileOp {
	return FileOp{
		Fs: afero.NewOsFs(),
	}
}

func (f *FileOp) GetFileSize(path string) (float64, error) {
	var size int64
	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return float64(size * 1024), nil

}

func GetParentMode(path string) (os.FileMode, error) {
	absPath, err := filepath.Abs(path)

	if err != nil {
		return 0, nil
	}

	for {
		info, err := os.Stat(absPath)
		if err == nil {
			return info.Mode() & os.ModePerm, nil
		}

		if !os.IsNotExist(err) {
			return 0, err
		}

		parentPath := filepath.Dir(absPath)

		if parentPath == absPath {
			return 0, fmt.Errorf("no existing directory found in the path: %s", path)
		}
		absPath = parentPath
	}

}
