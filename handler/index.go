package handler

import (
	dblayer "mycloud/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	openId, _ := c.Get("openId")
	//获取用户信息
	user := dblayer.GetUserInfo(openId)
	//获取用户仓库信息
	userFileStore := dblayer.GetUserFileStore(user.Id)
	//获取用户文件数量
	fileCount := dblayer.GetUserFileCount(user.FileStoreId)
	//获取用户文件夹数量
	fileFolderCount := dblayer.GetUserFileFolderCount(user.FileStoreId)
	//获取用户文件使用明细数量
	fileDetailUse := dblayer.GetFileDetailUse(user.FileStoreId)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"user":            user,
		"currIndex":       "active",
		"userFileStore":   userFileStore,
		"fileCount":       fileCount,
		"fileFolderCount": fileFolderCount,
		"fileDetailUse":   fileDetailUse,
	})
}
