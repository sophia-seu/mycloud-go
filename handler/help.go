package handler

import (
	dblayer "mycloud/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Help(c *gin.Context) {
	openId, _ := c.Get("openId")
	user := dblayer.GetUserInfo(openId)

	//获取用户文件使用明细数量
	fileDetailUse := dblayer.GetFileDetailUse(user.FileStoreId)

	c.HTML(http.StatusOK, "help.html", gin.H{
		"currHelp":      "active",
		"user":          user,
		"fileDetailUse": fileDetailUse,
	})
}
