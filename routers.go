package main

import (
	"GinProject12/controller/email"
	"GinProject12/controller/files"
	"GinProject12/controller/logs"
	"GinProject12/controller/ssh"
	systemtro "GinProject12/controller/system"
	"GinProject12/controller/user"
	_ "GinProject12/docs"
	"GinProject12/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(middleware.Cors())

	userGroup := r.Group("/api/user")
	{
		userGroup.POST("/register", user.BaseUserRegister)
		userGroup.POST("/login", user.BaseUserLogin)
		userGroup.GET("/info", middleware.AuthMiddleware(), user.BaseUserInfo)
	}

	adminGroup := r.Group("/api/admin")
	{
		adminGroup.POST("/register", user.AdminRegister)
		adminGroup.POST("/login", user.AdminLogin)
		adminGroup.GET("/info", middleware.AuthMiddleware(), user.AdminInfo)
	}

	sysGroup := r.Group("/sys")
	{
		sysGroup.GET("/info", systemtro.SysInfo)
		sysGroup.GET("/static", systemtro.SysStatic)
		sysGroup.GET("/net", systemtro.SysMonitorNet)
		sysGroup.GET("/io", systemtro.SysMonitorIo)
	}

	logGroup := r.Group("/logs")
	{
		logGroup.POST("/login", logs.LoginLog)
		logGroup.POST("/panel", logs.PanelLogin)
	}

	fileGroup := r.Group("/files")
	{
		fileGroup.POST("/search", filetro.FileDetail)
		fileGroup.POST("/content", filetro.GetFileContent)
		fileGroup.POST("/size", filetro.Size)
		fileGroup.POST("/rename", filetro.FileRename)
		fileGroup.POST("/recycle/search", filetro.RecycleBin)
		fileGroup.POST("/create", filetro.FileCreate)
		fileGroup.POST("/recover", filetro.FileRecover)
		fileGroup.POST("/remove", filetro.FileReMove)
		fileGroup.POST("/mode", filetro.FileChmod)
		fileGroup.POST("/compress", filetro.FilesCompress)
		fileGroup.POST("/decompress", filetro.FilesDecompress)
		fileGroup.GET("/download", filetro.Download)
		fileGroup.POST("/upload", filetro.UploadFiles)
	}

	emailGroup := r.Group("/email")
	{
		emailGroup.GET("/code", email.SendEmailCode)
		emailGroup.GET("/verify", email.VerifyCode)
	}

	sshGroup := r.Group("/ssh")
	{
		sshGroup.GET("/operate", sshtro.SshService)
		sshGroup.GET("/connect", sshtro.SshConnect)
	}

	return r
}
