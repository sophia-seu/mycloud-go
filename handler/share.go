package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	dblayer "mycloud/db"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lifei6671/gocaptcha"
)

//创建分享文件
func ShareFile(c *gin.Context) {
	openId, _ := c.Get("openId")
	//获取用户信息
	user := dblayer.GetUserInfo(openId)

	fId := c.Query("id")
	url := c.Query("url")
	//获取内容
	code := gocaptcha.RandText(4)

	fmt.Println(url)
	fileId, _ := strconv.Atoi(fId)
	hash := dblayer.CreateShare(code, user.UserName, fileId)

	c.JSON(http.StatusOK, gin.H{
		"url":  url + "?f=" + hash,
		"code": code,
	})
}

//分享文件页面
func SharePass(c *gin.Context) {
	f := c.Query("f")

	//获取分享信息
	shareInfo := dblayer.GetShareInfo(f)
	//获取文件信息
	file := dblayer.GetFileInfo(strconv.Itoa(shareInfo.FileId))

	c.HTML(http.StatusOK, "share.html", gin.H{
		"id":       shareInfo.FileId,
		"username": shareInfo.Username,
		"fileType": file.Type,
		"filename": file.FileName + file.Postfix,
		"hash":     shareInfo.Hash,
	})
}

//下载分享文件
func DownloadShareFile(c *gin.Context) {
	fileId := c.Query("id")
	code := c.Query("code")
	hash := c.Query("hash")

	fileInfo := dblayer.GetFileInfo(fileId)

	//校验提取码
	if ok := dblayer.VerifyShareCode(fileId, strings.ToLower(code)); !ok {
		c.Redirect(http.StatusMovedPermanently, "/file/share?f="+hash)
		return
	}

	//从本地拷贝数据到fileData
	fileData, err := ioutil.ReadFile(fileInfo.FilePath)
	if err != nil {
		log.Println("下载时读取文件失败...", err)
	}

	//下载次数+1
	dblayer.DownloadNumAdd(fileId)

	c.Header("Content-disposition", "attachment;filename=\""+fileInfo.FileName+fileInfo.Postfix+"\"")
	c.Data(http.StatusOK, "application/octect-stream", fileData)
}
