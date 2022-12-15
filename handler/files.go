package handler

import (
	"io/ioutil"
	"log"
	dblayer "mycloud/db"
	data "mycloud/meta"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Files(c *gin.Context) {

	// 1. 获取用户信息，拦截器中已经存入context内，可以直接从context读取
	openId, exist := c.Get("openId")
	if !exist {
		log.Println("Failed to get openId from context")
		return
	}
	fId := c.DefaultQuery("fId", "0")

	user := dblayer.GetUserInfo(openId)

	// 2. 获取该用户当前目录所有文件
	files := dblayer.GetUserFiles(fId, user.FileStoreId)

	// 3. 获取用户当前目录所有文件夹
	fileFolder := dblayer.GetFileFolders(fId, user.FileStoreId)

	// 4. 获取父级文件夹信息
	parentFolder := dblayer.GetParentFolder(fId)

	// 5. 获取当前目录所有父级
	currentAllParent := dblayer.GetCurrentAllParent(parentFolder, make([]data.FileFolder, 0))

	// 6. 获取当前目录信息
	currentFolder := dblayer.GetCurrentFolders(fId)

	// 7. 获取用户文件使用明细数量
	fileDetailUse := dblayer.GetFileDetailUse(user.FileStoreId)

	c.HTML(http.StatusOK, "files.html", gin.H{
		"currAll":          "active",
		"user":             user,
		"fId":              currentFolder.Id,
		"fName":            currentFolder.FileFolderName,
		"files":            files,
		"fileFolder":       fileFolder,
		"parentFolder":     parentFolder,
		"currentAllParent": currentAllParent,
		"fileDetailUse":    fileDetailUse,
	})

	// c.HTML(http.StatusOK, "files.html", nil)
}

func DocFiles(c *gin.Context) {

	openId, _ := c.Get("openId")
	user := dblayer.GetUserInfo(openId)

	//获取用户文件使用明细数量
	fileDetailUse := dblayer.GetFileDetailUse(user.FileStoreId)
	//获取文档类型文件
	docFiles := dblayer.GetTypeFile(1, user.FileStoreId)

	c.HTML(http.StatusOK, "doc-files.html", gin.H{
		"user":          user,
		"fileDetailUse": fileDetailUse,
		"docFiles":      docFiles,
		"docCount":      len(docFiles),
		"currDoc":       "active",
		"currClass":     "active",
	})
}

func ImageFiles(c *gin.Context) {
	openId, _ := c.Get("openId")
	user := dblayer.GetUserInfo(openId)

	//获取用户文件使用明细数量
	fileDetailUse := dblayer.GetFileDetailUse(user.FileStoreId)
	//获取图像类型文件
	imgFiles := dblayer.GetTypeFile(2, user.FileStoreId)

	c.HTML(http.StatusOK, "image-files.html", gin.H{
		"user":          user,
		"fileDetailUse": fileDetailUse,
		"imgFiles":      imgFiles,
		"imgCount":      len(imgFiles),
		"currImg":       "active",
		"currClass":     "active",
	})
}

func VideoFiles(c *gin.Context) {
	openId, _ := c.Get("openId")
	user := dblayer.GetUserInfo(openId)

	//获取用户文件使用明细数量
	fileDetailUse := dblayer.GetFileDetailUse(user.FileStoreId)
	//获取视频类型文件
	videoFiles := dblayer.GetTypeFile(3, user.FileStoreId)

	c.HTML(http.StatusOK, "video-files.html", gin.H{
		"user":          user,
		"fileDetailUse": fileDetailUse,
		"videoFiles":    videoFiles,
		"videoCount":    len(videoFiles),
		"currVideo":     "active",
		"currClass":     "active",
	})
}

func MusicFiles(c *gin.Context) {
	openId, _ := c.Get("openId")
	user := dblayer.GetUserInfo(openId)

	//获取用户文件使用明细数量
	fileDetailUse := dblayer.GetFileDetailUse(user.FileStoreId)
	//获取音频类型文件
	musicFiles := dblayer.GetTypeFile(4, user.FileStoreId)

	c.HTML(http.StatusOK, "music-files.html", gin.H{
		"user":          user,
		"fileDetailUse": fileDetailUse,
		"musicFiles":    musicFiles,
		"musicCount":    len(musicFiles),
		"currMusic":     "active",
		"currClass":     "active",
	})
}

func OtherFiles(c *gin.Context) {
	openId, _ := c.Get("openId")
	user := dblayer.GetUserInfo(openId)

	//获取用户文件使用明细数量
	fileDetailUse := dblayer.GetFileDetailUse(user.FileStoreId)
	//获取音频类型文件
	otherFiles := dblayer.GetTypeFile(5, user.FileStoreId)

	c.HTML(http.StatusOK, "other-files.html", gin.H{
		"user":          user,
		"fileDetailUse": fileDetailUse,
		"otherFiles":    otherFiles,
		"otherCount":    len(otherFiles),
		"currOther":     "active",
		"currClass":     "active",
	})
}

//处理新建文件夹
func AddFolder(c *gin.Context) {
	openId, _ := c.Get("openId")
	user := dblayer.GetUserInfo(openId)

	folderName := c.PostForm("fileFolderName")
	parentId := c.DefaultPostForm("parentFolderId", "0")

	//新建文件夹数据
	dblayer.CreateFolder(folderName, parentId, user.FileStoreId)

	//获取父文件夹信息
	parent := dblayer.GetParentFolder(parentId)

	c.Redirect(http.StatusMovedPermanently, "/cloud/files?fId="+parentId+"&fName="+parent.FileFolderName)
}

//修改文件夹名
func UpdateFileFolder(c *gin.Context) {
	fileFolderName := c.PostForm("fileFolderName")
	fileFolderId := c.PostForm("fileFolderId")

	fileFolder := dblayer.GetCurrentFolders(fileFolderId)

	dblayer.UpdateFolderName(fileFolderId, fileFolderName)

	c.Redirect(http.StatusMovedPermanently, "/cloud/files?fId="+strconv.Itoa(fileFolder.ParentFolderId))
}

//删除文件
func DeleteFile(c *gin.Context) {
	openId, _ := c.Get("openId")
	user := dblayer.GetUserInfo(openId)

	fId := c.DefaultQuery("fId", "")
	folderId := c.Query("folder")
	if fId == "" {
		return
	}

	//删除数据库文件数据
	dblayer.DeleteUserFile(fId, folderId, user.FileStoreId)

	c.Redirect(http.StatusMovedPermanently, "/cloud/files?fid="+folderId)
}

//删除文件夹
func DeleteFileFolder(c *gin.Context) {
	fId := c.DefaultQuery("fId", "")
	if fId == "" {
		return
	}
	//获取要删除的文件夹信息 取到父级目录重定向
	folderInfo := dblayer.GetCurrentFolders(fId)

	//删除文件夹并删除文件夹中的文件信息
	dblayer.DeleteFileFolder(fId)

	c.Redirect(http.StatusMovedPermanently, "/cloud/files?fId="+strconv.Itoa(folderInfo.ParentFolderId))
}

func DownloadFile(c *gin.Context) {
	fId := c.Query("fId")

	file := dblayer.GetFileInfo(fId)
	if file.FileHash == "" {
		return
	}

	//从oss获取文件
	// fileData := lib.DownloadOss(file.FileHash, file.Postfix)

	//从本地拷贝数据到fileData
	fileData, err := ioutil.ReadFile(file.FilePath)
	if err != nil {
		log.Println("下载时读取文件失败...", err)
	}
	//下载次数+1
	dblayer.DownloadNumAdd(fId)

	c.Header("Content-disposition", "attachment;filename=\""+file.FileName+file.Postfix+"\"")
	c.Data(http.StatusOK, "application/octect-stream", fileData)
}
