package files

import (
	"GinProject12/util"
	"GinProject12/util/cmd"
	"archive/zip"
	"fmt"
	"github.com/mholt/archiver/v4"
	"github.com/spf13/afero"
	"io"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type CompressType string

const (
	Zip      CompressType = "zip"
	Gz       CompressType = "gz"
	Bz2      CompressType = "bz2"
	Tar      CompressType = "tar"
	TarGz    CompressType = "tar.gz"
	Xz       CompressType = "xz"
	SdkZip   CompressType = "sdkZip"
	SdkTarGz CompressType = "sdkTarGz"
)

var (
	CONF ServerConfig
)

type ServerConfig struct {
	System    System    `mapstructure:"system"`
	LogConfig LogConfig `mapstructure:"log"`
}

type System struct {
	Port           string `mapstructure:"port"`
	Ipv6           string `mapstructure:"ipv6"`
	BindAddress    string `mapstructure:"bindAddress"`
	SSL            string `mapstructure:"ssl"`
	DbFile         string `mapstructure:"db_file"`
	DbPath         string `mapstructure:"db_path"`
	LogPath        string `mapstructure:"log_path"`
	DataDir        string `mapstructure:"data_dir"`
	TmpDir         string `mapstructure:"tmp_dir"`
	Cache          string `mapstructure:"cache"`
	Backup         string `mapstructure:"backup"`
	EncryptKey     string `mapstructure:"encrypt_key"`
	BaseDir        string `mapstructure:"base_dir"`
	Mode           string `mapstructure:"mode"`
	RepoUrl        string `mapstructure:"repo_url"`
	Version        string `mapstructure:"version"`
	Username       string `mapstructure:"username"`
	Password       string `mapstructure:"password"`
	Entrance       string `mapstructure:"entrance"`
	IsDemo         bool   `mapstructure:"is_demo"`
	AppRepo        string `mapstructure:"app_repo"`
	ChangeUserInfo string `mapstructure:"change_user_info"`
	OneDriveID     string `mapstructure:"one_drive_id"`
	OneDriveSc     string `mapstructure:"one_drive_sc"`
}

type LogConfig struct {
	Level     string `mapstructure:"level"`
	TimeZone  string `mapstructure:"timeZone"`
	LogName   string `mapstructure:"log_name"`
	LogSuffix string `mapstructure:"log_suffix"`
	MaxBackup int    `mapstructure:"max_backup"`
}

const (
	Deflate uint16 = 8 // DEFLATE compressed
)

type ZipArchiver struct{}

func NewZipArchiver() ShellArchiver {
	return &ZipArchiver{}
}

func (z ZipArchiver) Extract(filePath, dstDir string) error {
	if err := cmd.CheckCmdAvailability("unzip"); err != nil {
		return err
	}
	return cmd.ExecCmd(fmt.Sprintf("unzip -qo %s -d %s", filePath, dstDir))
}

func (z ZipArchiver) Compress(sourcePaths []string, dstFile string) error {
	var err error
	tmpFile := path.Join(CONF.System.TmpDir, fmt.Sprintf("%s%s.zip", util.RandStr(50), time.Now().Format("20060102150405")))
	op := afero.NewOsFs()
	defer func() {
		_ = op.Remove(tmpFile)
		if err != nil {
			_ = op.Remove(dstFile)
		}
	}()
	baseDir := path.Dir(sourcePaths[0])
	relativePaths := make([]string, len(sourcePaths))
	for i, sp := range sourcePaths {
		relativePaths[i] = path.Base(sp)
	}
	cmdStr := fmt.Sprintf("zip -qr %s  %s", tmpFile, strings.Join(relativePaths, " "))
	if err = cmd.ExecCmdWithDir(cmdStr, baseDir); err != nil {
		return err
	}
	if err = op.Rename(tmpFile, dstFile); err != nil {
		return err
	}
	return nil
}

func ZipFile(files []archiver.File, dst afero.File) error {

	zw := zip.NewWriter(dst)
	defer zw.Close()

	for _, file := range files {
		hdr, err := zip.FileInfoHeader(file)
		if err != nil {
			return err
		}
		hdr.Method = zip.Deflate
		hdr.Name = file.NameInArchive
		// 文件夹
		if file.IsDir() {
			if !strings.HasSuffix(hdr.Name, "/") {
				hdr.Name += "/"
			}
		}
		w, err := zw.CreateHeader(hdr)
		if err != nil {
			return err
		}
		if file.IsDir() {
			continue
		}
		// 链接
		if file.LinkTarget != "" {
			_, err = w.Write([]byte(filepath.ToSlash(file.LinkTarget)))
			if err != nil {
				return err
			}
		} else {
			// 普通文件
			fileReader, err := file.Open()
			if err != nil {
				return err
			}
			_, err = io.Copy(w, fileReader)
			if err != nil {
				return err
			}
		}
	}
	return nil

}
