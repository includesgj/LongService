package filetro

import (
	"GinProject12/model"
	"GinProject12/response"
	files "GinProject12/serverce/cmd/files"
	"GinProject12/util"
	futil "GinProject12/util/files"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/afero"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
)

var (
	service files.FileService
)

// FileDetail 获取该路径的文件
// @Summary      获取该路径的文件
// @Description  获取该路径的文件
// @Tags         File
// @Accept       json
// @Produce      json
// @Param		 FileOption body files.FileOption ture "request"
// @Success      200  {object} files.FileInfo
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /files/search [POST]
func FileDetail(c *gin.Context) {

	req := files.FileOption{}

	if err := util.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	list, err := service.GetFileList(req)

	if err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, err.Error())
		return
	}

	response.Success(c, gin.H{"data": list}, "成功")

}

// GetFileContent 查看文件属性和查看文件内容
// @Summary      查看文件属性和查看文件内容
// @Description  查看文件属性和查看文件内容
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        content body files.FileContentReq ture "request"
// @Success      200  {object}  files.FileInfo
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /files/content [POST]
func GetFileContent(c *gin.Context) {
	req := files.FileContentReq{}
	if err := util.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	content, err := service.GetContent(req)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}

	response.Success(c, gin.H{"data": content}, "成功")

}

// Size 查看文件大小
// @Summary      查看文件大小
// @Description  查看文件大小
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        req body files.DirSizeReq ture "request"
// @Success      200  {object}  files.DirSizeRes
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /files/size [POST]
func Size(c *gin.Context) {
	req := files.DirSizeReq{}
	if err := util.CheckBindAndValidate(&req, c); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, err.Error())
		return
	}
	res, err := service.GetDirSize(req)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}
	response.Success(c, gin.H{"data": res}, "成功")
}

// FileRename 修改文件名字
// @Summary      修改文件名字
// @Description  修改文件名字
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        req body files.FileRenameReq ture "request"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /files/rename [POST]
func FileRename(c *gin.Context) {
	req := files.FileRenameReq{}
	if err := util.CheckBindAndValidate(&req, c); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, err.Error())
		return
	}

	if err := service.Rename(req); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}

	response.Success(c, nil, "成功")

}

// FileReMove 删除文件
// @Summary      删除文件
// @Description  删除文件
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        req body files.RemoveReq ture "request"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /files/remove [POST]
func FileReMove(c *gin.Context) {
	req := files.RemoveReq{}
	if err := util.CheckBindAndValidate(&req, c); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, err.Error())
		return
	}

	if err := service.Remove(req); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}
	response.Success(c, nil, "成功")
}

// FileCreate 创建文件
// @Summary      创建文件
// @Description  创建文件 IsSymLink软链接 IsLink硬链接
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        req body files.CreateReq ture "request"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /files/create [POST]
func FileCreate(c *gin.Context) {
	req := files.CreateReq{}
	if err := util.CheckBindAndValidate(&req, c); err != nil {
		return
	}
	if err := service.Create(req); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}
	response.Success(c, nil, "成功")

}

// FileRecover 恢复回收站(单个文件)
// @Summary      恢复回收站(单个文件)
// @Description  恢复回收站(单个文件)
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        req body model.RecoverReq ture "request"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /files/recover [POST]
func FileRecover(c *gin.Context) {
	req := model.RecoverReq{}
	if err := util.CheckBindAndValidate(&req, c); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, err.Error())
		return
	}

	if err := service.RecoverInfo(req); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}
	response.Success(c, nil, "成功")

}

// RecycleBin 查询回收站
// @Summary      查询回收站(分页)
// @Description  查询回收站(分页)
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        req body model.PageInfo ture "request"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /files/recycle/search [POST]
func RecycleBin(c *gin.Context) {
	req := model.PageInfo{}
	if err := util.CheckBindAndValidate(&req, c); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, err.Error())
		return
	}

	list, err := service.GetRecycleBin(req)

	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}
	response.Success(c, gin.H{"data": list}, "成功")
}

// FileChmod 修改文件权限
// @Summary      修改文件权限
// @Description  修改文件权限
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        req body files.ChmodReq ture "request"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /files/mode [POST]
func FileChmod(c *gin.Context) {
	req := files.ChmodReq{}
	if err := util.CheckBindAndValidate(&req, c); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, err.Error())
		return
	}

	if err := service.Chmod(req); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}

	response.Success(c, nil, "成功")
}

// FilesCompress 压缩文件
// @Summary      压缩文件
// @Description  压缩文件
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        req body files.FileCompress ture "request"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /files/compress [POST]
func FilesCompress(c *gin.Context) {
	req := files.FileCompress{}
	if err := util.CheckBindAndValidate(&req, c); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, err.Error())
		return
	}

	if err := service.Compress(req); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}
	response.Success(c, nil, "成功")

}

// FilesDecompress 解压文件
// @Summary      解压文件
// @Description  解压文件
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        req body files.FileDecompress ture "request"
// @Success      200
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /files/decompress [POST]
func FilesDecompress(c *gin.Context) {
	req := files.FileDecompress{}
	if err := util.CheckBindAndValidate(&req, c); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, err.Error())
		return
	}
	if err := service.Decompress(req); err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}
	response.Success(c, nil, "成功")
}

// Download 下载文件
// @Tags File
// @Summary      下载文件
// @Description  下载文件指明路径
// @Accept json
// @Param  		 path query string ture "path"
// @Success 200
// @Router /files/download [get]
// @Router       /files/download [GET]
func Download(c *gin.Context) {
	filePath := c.Query("path")
	file, err := os.Open(filePath)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}
	defer file.Close()
	info, _ := file.Stat()
	c.Header("Content-Length", strconv.FormatInt(info.Size(), 10))
	c.Header("Content-Disposition", "attachment; filename*=utf-8''"+url.PathEscape(info.Name()))
	http.ServeContent(c.Writer, c.Request, info.Name(), info.ModTime(), file)
}

// UploadFiles 上传文件
// @Tags File
// @Summary 上传文件
// @Description 上传文件指明路径 file 上传的文件 path存放的路径 cover 是否覆盖 用 postman的 body form-data
// @Param file formData file true "request"
// @Success 200
// @Router /files/upload [post]
func UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()

	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}

	upFiles := form.File["file"]
	upPath := form.Value["path"]

	cover := true
	if cv, ok := form.Value["cover"]; !ok {
		if len(cv) != 0 {
			parseBool, _ := strconv.ParseBool(cv[0])
			cover = parseBool
		}
	}

	if len(upPath) == 0 || !strings.Contains(upPath[0], "/") {
		response.Response(c, http.StatusInternalServerError, 400, nil, errors.New("路径错误").Error())
		return
	}
	dir := path.Dir(upPath[0])

	_, err = os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		mode, err := futil.GetParentMode(dir)
		if err != nil {
			response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
			return
		}
		if err = os.MkdirAll(dir, mode); err != nil {
			response.Response(c, http.StatusInternalServerError, 400, nil, fmt.Errorf("创建文件夹 %s 失败: %v", dir, err).Error())
			return
		}
	}
	info, err := os.Stat(dir)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, err.Error())
		return
	}
	mode := info.Mode()

	fileOp := afero.NewOsFs()

	success := 0
	failures := make(map[string]error)
	for _, file := range upFiles {
		dstFilename := path.Join(upPath[0], file.Filename)
		dstDir := path.Dir(dstFilename)
		if _, err = fileOp.Stat(dstDir); err != nil {
			if err = fileOp.MkdirAll(dstDir, mode); err != nil {
				e := fmt.Errorf("创建文件 [%s] 失败: %v", path.Dir(dstFilename), err)
				failures[file.Filename] = e
				log.Println(e)
				continue
			}
		}
		tmpFilename := dstFilename + ".tmp"
		if err := c.SaveUploadedFile(file, tmpFilename); err != nil {
			_ = os.Remove(tmpFilename)
			e := fmt.Errorf("上传 [%s] 文件失败:%v", file.Filename, err)
			failures[file.Filename] = e
			log.Println(e)
			continue
		}
		dstInfo, statErr := os.Stat(dstFilename)
		if cover {
			_ = os.Remove(dstFilename)
		}

		err = os.Rename(tmpFilename, dstFilename)
		if err != nil {
			_ = os.Remove(tmpFilename)
			e := fmt.Errorf("上传 [%s] 文件失败: %v", file.Filename, err)
			failures[file.Filename] = e
			log.Println(e)
			continue
		}
		if statErr == nil {
			_ = os.Chmod(dstFilename, dstInfo.Mode())
		} else {
			_ = os.Chmod(dstFilename, mode)
		}
		success++
	}
	if success == 0 {
		response.Response(c, http.StatusInternalServerError, 500, nil, failures)
	} else {
		response.Success(c, nil, fmt.Sprintf("%d 文件上传成功", success))
	}

}
