package files

import (
	"github.com/spf13/afero"
	"io/fs"
	"os"
	"time"
)

var (
	ErrFileCanNotRead = "ErrFileCanNotRead"
)

const dotCharacter = 46

type FileInfo struct {
	Fs         afero.Fs    `json:"-"`
	Path       string      `json:"path"`
	Name       string      `json:"name"`
	User       string      `json:"user"`
	Group      string      `json:"group"`
	Uid        string      `json:"uid"`
	Gid        string      `json:"gid"`
	Extension  string      `json:"extension"`
	Content    string      `json:"content"`
	Size       int64       `json:"size"`
	IsDir      bool        `json:"isDir"`
	IsSymlink  bool        `json:"isSymlink"`
	IsHidden   bool        `json:"isHidden"`
	LinkPath   string      `json:"linkPath"`
	Type       string      `json:"type"`
	Mode       string      `json:"mode"`
	MimeType   string      `json:"mimeType"`
	UpdateTime time.Time   `json:"updateTime"`
	ModTime    time.Time   `json:"modTime"`
	FileMode   os.FileMode `json:"-"`
	Items      []*FileInfo `json:"items"`
	ItemTotal  int         `json:"itemTotal"`
	FavoriteID uint        `json:"favoriteID"`
}

type FileService struct{}

/*
Path string：指定的文件路径。
Search string：搜索关键词。
ContainSub bool：是否包含子文件夹。
Expand bool：是否展开文件夹。
Dir bool：是否只显示文件夹。
ShowHidden bool：是否显示隐藏文件。
Page int：请求的页数。
PageSize int：每页显示的文件数量。
SortBy string：按照什么方式排序文件列表。
SortOrder string：排序顺序（升序或降序）。
*/

type FileOption struct {
	Path       string `json:"path" validate:"required"`
	Search     string `json:"search"`
	ContainSub bool   `json:"containSub"`
	Expand     bool   `json:"expand"`
	Dir        bool   `json:"dir"`
	ShowHidden bool   `json:"showHidden"`
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
	SortBy     string `json:"sortBy"`
	SortOrder  string `json:"sortOrder"`
}

type FileSearchInfo struct {
	Path string `json:"path"`
	fs.FileInfo
}

type FileContentReq struct {
	Path string `json:"path" validate:"required"`
}

func (f *FileService) GetFileList(op FileOption) (FileInfo, error) {
	var fileInfo FileInfo
	if _, err := os.Stat(op.Path); err != nil && os.IsNotExist(err) {
		return fileInfo, nil
	}
	info, err := NewFileInfo(op)
	if err != nil {
		return fileInfo, err
	}
	fileInfo = *info
	return fileInfo, nil
}

func (f *FileService) GetContent(op FileContentReq) (FileInfo, error) {
	info, err := NewFileInfo(FileOption{
		Path:   op.Path,
		Expand: true,
	})
	if err != nil {
		return FileInfo{}, err
	}
	return *info, nil
}

type DirSizeReq struct {
	Path string `json:"path" validate:"required"`
}

type DirSizeRes struct {
	Size float64 `json:"size"`
}

type FileRenameReq struct {
	OldName string `json:"oldName" validate:"required"`
	NewName string `json:"newName" validate:"required"`
	Path    string `json:"path" validate:"required"`
}

type RemoveReq struct {
	Path    string `json:"path" validate:"required"`
	RealDel bool   `json:"realDel"`
	IsDir   bool   `json:"isDir"`
}

type CreateReq struct {
	Path      string `json:"path" validate:"required"`
	Mode      int    `json:"mode"` // 十进制
	IsLink    bool   `json:"isLink"`
	IsDir     bool   `json:"isDir"`
	IsSymLink bool   `json:"isSymLink"`
	LinkPath  string `json:"linkPath"`
}

type ChmodReq struct {
	Path string `json:"path" validate:"required"`
	Mode int    `json:"mode" validate:"required, max=511"` // 十进制
}

type FileCompress struct {
	Dst     string   `json:"dst" validate:"required"`
	Name    string   `json:"name"  validate:"required"`
	Type    string   `json:"type"  validate:"required"`
	Files   []string `json:"files"  validate:"required"`
	Replace bool     `json:"replace"`
}

type FileDecompress struct {
	Type string `json:"type" validate:"required"`
	Dst  string `json:"dst" validate:"required"`
	Path string `json:"path" validate:"required"`
}
