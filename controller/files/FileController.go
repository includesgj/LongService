package filetro

import (
	"GinProject12/model"
	"GinProject12/response"
	files "GinProject12/serverce/cmd/files"
	"GinProject12/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	service files.FileService
)

// FileDetail 获取该路径的文件
// @Summary      获取该路径的文件
// @Description  获取该路径的文件
// @Tags         file
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

// GetFileContent 查看文件属性
// @Summary      查看文件属性
// @Description  查看文件属性
// @Tags         file
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
// @Tags         file
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
// @Tags         file
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
// @Tags         file
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
// @Description  创建文件
// @Tags         file
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
		response.Response(c, http.StatusBadRequest, 400, nil, err.Error())
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
// @Tags         file
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
// @Tags         file
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
// @Tags         file
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
// @Tags         file
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
// @Tags         file
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
