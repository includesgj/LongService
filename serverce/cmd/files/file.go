package files

import (
	sdb "GinProject12/databases"
	"GinProject12/model"
	"GinProject12/util"
	"GinProject12/util/files"
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/mholt/archiver/v4"
	_ "github.com/mholt/archiver/v4"
	"github.com/spf13/afero"
	_ "github.com/spf13/afero"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func (f *FileInfo) search(search string, count int) (files []FileSearchInfo, total int, err error) {
	cmd := exec.Command("find", f.Path, "-name", fmt.Sprintf("*%s*", search))
	output, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	if err = cmd.Start(); err != nil {
		return
	}
	defer func() {
		_ = cmd.Wait()
		_ = cmd.Process.Kill()
	}()

	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		line := scanner.Text()
		info, err := os.Stat(line)
		if err != nil {
			continue
		}
		total++
		if total > count {
			continue
		}
		files = append(files, FileSearchInfo{
			Path:     line,
			FileInfo: info,
		})
	}
	if err = scanner.Err(); err != nil {
		return
	}
	return
}

func (f *FileInfo) listChildren(option FileOption) error {
	afs := &afero.Afero{Fs: f.Fs}
	var (
		files []FileSearchInfo
		err   error
		total int
	)

	if option.Search != "" && option.ContainSub {
		files, total, err = f.search(option.Search, option.Page*option.PageSize)
		if err != nil {
			return err
		}
	} else {
		dirFiles, err := afs.ReadDir(f.Path)
		if err != nil {
			return err
		}
		var (
			dirs     []FileSearchInfo
			fileList []FileSearchInfo
		)
		for _, file := range dirFiles {
			info := FileSearchInfo{
				Path:     f.Path,
				FileInfo: file,
			}
			if file.IsDir() {
				dirs = append(dirs, info)
			} else {
				fileList = append(fileList, info)
			}
		}
		sortFileList(dirs, option.SortBy, option.SortOrder)
		sortFileList(fileList, option.SortBy, option.SortOrder)
		files = append(dirs, fileList...)
	}

	var items []*FileInfo
	for _, df := range files {
		if option.Dir && !df.IsDir() {
			continue
		}
		name := df.Name()
		fPath := path.Join(df.Path, df.Name())
		if option.Search != "" {
			if option.ContainSub {
				fPath = df.Path
				name = strings.TrimPrefix(strings.TrimPrefix(fPath, f.Path), "/")
			} else {
				lowerName := strings.ToLower(name)
				lowerSearch := strings.ToLower(option.Search)
				if !strings.Contains(lowerName, lowerSearch) {
					continue
				}
			}
		}
		if !option.ShowHidden && IsHidden(name) {
			continue
		}
		f.ItemTotal++
		isSymlink, isInvalidLink := false, false
		if IsSymlink(df.Mode()) {
			isSymlink = true
			info, err := f.Fs.Stat(fPath)
			if err == nil {
				df.FileInfo = info
			} else {
				isInvalidLink = true
			}
		}

		file := &FileInfo{
			Fs:        f.Fs,
			Name:      name,
			Size:      df.Size(),
			ModTime:   df.ModTime(),
			FileMode:  df.Mode(),
			IsDir:     df.IsDir(),
			IsSymlink: isSymlink,
			IsHidden:  IsHidden(fPath),
			Extension: filepath.Ext(name),
			Path:      fPath,
			Mode:      fmt.Sprintf("%04o", df.Mode().Perm()),
			User:      GetUsername(df.Sys().(*syscall.Stat_t).Uid),
			Group:     GetGroup(df.Sys().(*syscall.Stat_t).Gid),
			Uid:       strconv.FormatUint(uint64(df.Sys().(*syscall.Stat_t).Uid), 10),
			Gid:       strconv.FormatUint(uint64(df.Sys().(*syscall.Stat_t).Gid), 10),
		}

		if isSymlink {
			file.LinkPath = GetSymlink(fPath)
		}
		if df.Size() > 0 {
			file.MimeType = GetMimeType(fPath)
		}
		if isInvalidLink {
			file.Type = "invalid_link"
		}
		items = append(items, file)
	}
	if option.ContainSub {
		f.ItemTotal = total
	}
	start := (option.Page - 1) * option.PageSize
	end := option.PageSize + start
	var result []*FileInfo
	if start < 0 || start > f.ItemTotal || end < 0 || start > end {
		result = items
	} else {
		if end > f.ItemTotal {
			result = items[start:]
		} else {
			result = items[start:end]
		}
	}

	f.Items = result
	return nil
}

func (f *FileInfo) getContent() error {
	if IsBlockDevice(f.FileMode) {
		return errors.New(ErrFileCanNotRead)
	}
	if f.Size > 1024*1024*10 {
		return errors.New("ErrFileToLarge")
	}
	afs := &afero.Afero{Fs: f.Fs}
	cByte, err := afs.ReadFile(f.Path)
	if err != nil {
		return nil
	}
	if len(cByte) > 0 && DetectBinary(cByte) {
		return errors.New(ErrFileCanNotRead)
	}
	f.Content = string(cByte)
	return nil
}

// GetDirSize 获取文件夹的大小Byte单位
func (f *FileService) GetDirSize(req DirSizeReq) (DirSizeRes, error) {

	res := DirSizeRes{}

	// 虚拟目录
	if req.Path == "/proc" {
		return res, nil
	}

	// du -s 可查看当前文件夹的大小
	command := exec.Command("du", "-s", req.Path)

	output, err := command.Output()

	if err == nil {
		fields := strings.Fields(string(output))
		if len(fields) == 2 {
			var size int
			_, err = fmt.Sscanf(fields[0], "%d", &size)
			if err == nil {
				// 为与GetDirSize 单位一致为Byte 所以乘1024
				res.Size = float64(size * 1024)
				return res, nil
			}
		}
	}

	op := futil.NewFileOp()
	size, err := op.GetFileSize(req.Path)
	if err != nil {
		return res, err
	}

	res.Size = size

	return res, err

}

// Rename 修改文件名称
func (f *FileService) Rename(req FileRenameReq) error {
	oldPath := req.Path + req.OldName
	newPath := req.Path + req.NewName
	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}
	return nil
}

// Remove
/**
-- 前端传过来的
{
	path 文件路径
	realDel 是否加入回收站 false为放入回收站
	isDir 是否文件夹
}
把需要到回收站的放到 根目录下的.soul_power 文件命名 (删除时间 原来的路径 别忘了名称也算 分开写 )
*/
func (f *FileService) Remove(req RemoveReq) error {
	// 真删
	if req.RealDel {
		if err := os.Remove(req.Path); err != nil {
			return err
		}
	}

	// 假删
	// 先创建存放回收站的文件夹 .soul_power
	if err := FileIsExist("/.soul_power"); err == nil {
		if err = CreateDir(CreateReq{Path: "./soul_power", Mode: 0755}); err != nil {
			return err
		}
	}

	paths := strings.Split(req.Path, "/")
	rName := strings.Join(paths, "_sp_")

	deleteTime := time.Now()

	op := futil.NewFileOp()

	size, err := op.GetFileSize(req.Path)
	if err != nil {
		return err
	}

	newPath := fmt.Sprintf("_sp_%s%s_p_%d_%d", "file", rName, int(size), deleteTime.Unix())
	if err = os.Rename(req.Path, newPath); err != nil {
		return err
	}
	return nil
}

// Create 创建文件
func (f *FileService) Create(req CreateReq) error {

	if err := FileIsExist(req.Path); err != nil {
		return err
	}

	if req.Mode == 0 {
		// 如果有父就继承父
		fileInfo, err := os.Stat(filepath.Dir(req.Path))
		if err == nil {
			req.Mode = int(fileInfo.Mode().Perm())
		} else {
			req.Mode = 0755
		}
	}

	// 此时没有同名的可以创建
	if req.IsDir {
		if err := CreateDir(req); err != nil {
			return err
		}
	} else if req.IsLink {
		if err := CreateLink(req); err != nil {
			return err
		}
	} else {
		if err := CreateFile(req); err != nil {
			return err
		}
	}

	return nil
}

// Chmod 修改文件权限
func (f *FileService) Chmod(req ChmodReq) error {
	mode, err := util.TransformOctal(req.Mode)
	if err != nil {
		return err
	}

	err = os.Chmod(req.Path, mode)
	if err != nil {
		return err
	}
	return nil
}

// GetRecycleBin 获取回收站的所有数据(按分页)
func (f *FileService) GetRecycleBin(info model.PageInfo) ([]model.RecycleBin, error) {
	total, err := sdb.RecycleBinPage(info)
	if err != nil {
		return nil, err
	}

	return total, nil
}

// RecoverInfo 恢复回收站
func (f *FileService) RecoverInfo(req model.RecoverReq) error {
	// 查询数据库查找到指定的然后恢复然后删除这一条
	info, err := sdb.RecycleBinInfo(req)
	if err != nil {
		return err
	}

	if err = os.Rename(path.Join(info.From, info.RName), info.SourcePath); err != nil {
		return err
	}

	err = sdb.DelRecycleBin(info.Id)

	if err != nil {
		return err
	}

	return nil

}

// Compress 压缩文件
func (f *FileService) Compress(s FileCompress) error {
	fs := afero.NewOsFs()

	_, err := fs.Stat(filepath.Join(s.Dst, s.Name))

	if err != nil && !s.Replace {
		return err
	}
	cType := Zip

	format := getFormat(cType)

	baseNameMap := make(map[string]string, len(s.Files))
	for _, path := range s.Files {
		Base := filepath.Base(path)
		baseNameMap[path] = Base
	}

	if _, err = fs.Stat(s.Dst); err != nil {
		_, _ = fs.Create(s.Dst)
	}

	files, err := archiver.FilesFromDisk(nil, baseNameMap)
	if err != nil {
		return err
	}

	dstPath := filepath.Join(s.Dst, s.Name)
	out, err := fs.Create(dstPath)
	if err != nil {
		return err
	}

	switch cType {
	case Zip:
		if err := ZipFile(files, out); err == nil {
			return err
		}
		_ = fs.Remove(s.Dst)
		op := ZipArchiver{}
		op.Compress(s.Files, dstPath)
		return err
	default:
		err = format.Archive(context.Background(), out, files)
		if err != nil {
			_ = fs.Remove(s.Dst)
			return err
		}
	}
	return nil

}

// Decompress 解压文件
func (f *FileService) Decompress(s FileDecompress) error {
	cType := CompressType(s.Type)
	if err := decompressWithSDK(s.Path, s.Dst, cType); err != nil {
		if cType == Tar || cType == Zip {
			shellArchiver, err := NewShellArchiver(cType)
			if err != nil {
				return err
			}
			return shellArchiver.Extract(s.Path, s.Dst)
		}
		return err
	}
	return nil
}
