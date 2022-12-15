package router

import (
	"mycloud/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoute() *gin.Engine {

	router := gin.Default()

	router.GET("/cloud/register", handler.Register)
	router.POST("/cloud/register", handler.HandlerRegister)

	router.GET("/cloud/login", handler.Login)
	router.POST("/cloud/login", handler.HandlerLogin)

	router.GET("/file/share", handler.SharePass)
	router.GET("/file/shareDownload", handler.DownloadShareFile)

	cloud := router.Group("cloud")

	cloud.Use(handler.CheckLogin)
	{
		cloud.GET("/index", handler.Index)
		cloud.GET("/files", handler.Files)
		cloud.GET("/upload", handler.Upload)
		cloud.GET("/doc-files", handler.DocFiles)
		cloud.GET("/image-files", handler.ImageFiles)
		cloud.GET("/video-files", handler.VideoFiles)
		cloud.GET("/music-files", handler.MusicFiles)
		cloud.GET("/other-files", handler.OtherFiles)
		cloud.GET("/deleteFile", handler.DeleteFile)
		cloud.GET("/deleteFolder", handler.DeleteFileFolder)
		cloud.GET("/downloadFile", handler.DownloadFile)
		cloud.GET("/logout", handler.Logout)
		// cloud.GET("/help", handler.Help)
	}

	{
		cloud.POST("/uploadFile", handler.HandlerUpload)
		cloud.POST("/addFolder", handler.AddFolder)
		cloud.POST("/updateFolder", handler.UpdateFileFolder)
		cloud.POST("/getQrCode", handler.ShareFile)
	}

	return router
}
