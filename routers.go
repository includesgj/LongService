package main

import (
	"GinProject12/controller/email"
	"GinProject12/controller/files"
	"GinProject12/controller/logs"
	monitortro "GinProject12/controller/monitor"
	patroltro "GinProject12/controller/patrol"
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
		sysGroup.GET("/info", middleware.AuthMiddleware(), systemtro.SysInfo)
		sysGroup.GET("/static", middleware.AuthMiddleware(), systemtro.SysStatic)
		sysGroup.GET("/net", middleware.AuthMiddleware(), systemtro.SysMonitorNet)
		sysGroup.GET("/io", middleware.AuthMiddleware(), systemtro.SysMonitorIo)
		sysGroup.GET("/name", middleware.AuthMiddleware(), systemtro.GetNetOrDiskName)
	}

	logGroup := r.Group("/logs")
	{
		logGroup.POST("/login", middleware.AuthMiddleware(), logs.LoginLog)
		logGroup.POST("/panel", middleware.AuthMiddleware(), logs.PanelLogin)
	}

	fileGroup := r.Group("/files")
	{
		fileGroup.POST("/search", middleware.AuthMiddleware(), filetro.FileDetail)
		fileGroup.POST("/content", middleware.AuthMiddleware(), filetro.GetFileContent)
		fileGroup.POST("/size", middleware.AuthMiddleware(), filetro.Size)
		fileGroup.POST("/rename", middleware.AdminAuthMiddleware(), filetro.FileRename)
		fileGroup.POST("/recycle/search", middleware.AuthMiddleware(), filetro.RecycleBin)
		fileGroup.POST("/create", middleware.AdminAuthMiddleware(), filetro.FileCreate)
		fileGroup.POST("/recover", middleware.AdminAuthMiddleware(), filetro.FileRecover)
		fileGroup.POST("/remove", middleware.AdminAuthMiddleware(), filetro.FileReMove)
		fileGroup.POST("/mode", middleware.AdminAuthMiddleware(), filetro.FileChmod)
		fileGroup.POST("/compress", middleware.AdminAuthMiddleware(), filetro.FilesCompress)
		fileGroup.POST("/decompress", middleware.AdminAuthMiddleware(), filetro.FilesDecompress)
		fileGroup.GET("/download", middleware.AuthMiddleware(), filetro.Download)
		fileGroup.POST("/upload", middleware.AuthMiddleware(), filetro.UploadFiles)
	}

	emailGroup := r.Group("/email")
	{
		emailGroup.GET("/code", email.SendEmailCode)
		emailGroup.GET("/verify", email.VerifyCode)
	}

	sshGroup := r.Group("/ssh")
	{
		sshGroup.GET("/operate", middleware.AdminAuthMiddleware(), sshtro.SshService)
		sshGroup.GET("/connect", middleware.AdminAuthMiddleware(), sshtro.SshConnect)
	}
	monitorGroup := r.Group("/monitor")
	{
		monitorGroup.POST("/add", middleware.AuthMiddleware(), monitortro.Monitor)
		monitorGroup.GET("/del", middleware.AuthMiddleware(), monitortro.DelMonitor)
		monitorGroup.POST("/sel", middleware.AuthMiddleware(), monitortro.SelectMonitor)
	}
	patrolGroup := r.Group("/patrol")
	{
		monitorGroup.POST("/add", middleware.AdminAuthMiddleware(), patroltro.AddPatrol)
		monitorGroup.GET("/del", middleware.AdminAuthMiddleware(), patroltro.DelPatrol)
		monitorGroup.POST("/sel", middleware.AuthMiddleware(), patroltro.PatrolPage)
		patrolGroup.Group("/user")
		{
			monitorGroup.POST("/add", middleware.AuthMiddleware(), patroltro.AddPatrolInfo)
			monitorGroup.GET("/del", middleware.AuthMiddleware(), patroltro.DelPatrolInfo)
			monitorGroup.POST("/sel", middleware.AuthMiddleware(), patroltro.PatrolInfoPage)
		}
	}

	return r
}
