package handler

import (
	"fmt"
	"io"
	dblayer "mycloud/db"
	"mycloud/lib"
	data "mycloud/meta"
	"mycloud/util"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

//上传文件页面
func Upload(c *gin.Context) {
	openId, _ := c.Get("openId")
	fId := c.DefaultQuery("fId", "0")
	//获取用户信息
	user := dblayer.GetUserInfo(openId)
	//获取当前目录信息
	currentFolder := dblayer.GetCurrentFolders(fId)
	//获取当前目录所有的文件夹信息
	fileFolders := dblayer.GetFileFolders(fId, user.FileStoreId)
	//获取父级的文件夹信息
	parentFolder := dblayer.GetParentFolder(fId)
	//获取当前目录所有父级
	currentAllParent := dblayer.GetCurrentAllParent(parentFolder, make([]data.FileFolder, 0))
	//获取用户文件使用明细数量
	fileDetailUse := dblayer.GetFileDetailUse(user.FileStoreId)

	c.HTML(http.StatusOK, "upload.html", gin.H{
		"user":             user,
		"currUpload":       "active",
		"fId":              currentFolder.Id,
		"fName":            currentFolder.FileFolderName,
		"fileFolders":      fileFolders,
		"parentFolder":     parentFolder,
		"currentAllParent": currentAllParent,
		"fileDetailUse":    fileDetailUse,
	})
}

func HandlerUpload(c *gin.Context) {
	openId, _ := c.Get("openId")

	// 1. 获取用户信息
	user := dblayer.GetUserInfo(openId)

	Fid := c.GetHeader("id")
	conf := lib.LoadServerConfig()

	//接收上传文件
	file, head, err := c.Request.FormFile("file")

	if err != nil {
		fmt.Println("文件上传错误", err.Error())
		return
	}
	defer file.Close()

	//判断当前文件夹是否有同名文件
	if ok := dblayer.CurrFileExists(Fid, head.Filename); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": 501,
		})
		return
	}

	//判断用户的容量是否足够
	if ok := dblayer.CapacityIsEnough(head.Size, user.FileStoreId); !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": 503,
		})
		return
	}

	//文件保存本地的路径
	user_dir := filepath.Join(conf.UploadLocation, user.UserName)

	if is_exist, _ := util.PathExists(user_dir); !is_exist {
		os.Mkdir(user_dir, os.ModePerm)
	}

	location := filepath.ToSlash(filepath.Join(user_dir, head.Filename))

	//在本地创建一个新的文件
	newFile, err := os.Create(location)
	if err != nil {
		fmt.Println("文件创建失败", err.Error())
		return
	}
	defer newFile.Close()

	//将上传文件拷贝至新创建的文件中
	fileSize, err := io.Copy(newFile, file)
	if err != nil {
		fmt.Println("文件拷贝错误", err.Error())
		return
	}

	//游标回到文件的开头，计算文件哈希值
	_, _ = newFile.Seek(0, 0)
	fileHash := util.GetSHA256HashCode(newFile)

	// //写入ceph存储
	// newFile.Seek(0, 0)
	// copy_data, _ := ioutil.ReadAll(newFile)
	// // 该文件存入的位置
	// cepg_store_location := filepath.ToSlash(filepath.Join("ceph", head.Filename))
	// //存储到ceph路径中
	// _ = lib.PutObject(user.UserName, cepg_store_location, copy_data)
	// // TOFO 接下来应该修改ceph路径到文件表中，但是文件表目前没有关于ceph的内容，后续再完善。。。

	//新建文件信息
	dblayer.CreateFile(head.Filename, fileHash, fileSize, Fid, user.FileStoreId, location)
	//上传成功减去相应剩余容量
	dblayer.SubtractSize(fileSize/1024, user.FileStoreId)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
	})

}
