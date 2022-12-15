package db

import (
	"fmt"
	"mycloud/db/mysql"
	data "mycloud/meta"
	"mycloud/util"
	"path"
	"strconv"
	"strings"
	"time"
)

//获取父类的id
func GetParentFolder(fId string) (fileFolder data.FileFolder) {
	mysql.DB.Find(&fileFolder, "id = ?", fId)
	return
}

//获取用户文件夹数量
func GetUserFileFolderCount(fileStoreId int) (fileFolderCount int) {
	var fileFolder []data.FileFolder
	mysql.DB.Find(&fileFolder, "file_store_id = ?", fileStoreId).Count(&fileFolderCount)
	return
}

//获取目录所有文件夹
func GetFileFolders(parentId string, fileStoreId int) (fileFolders []data.FileFolder) {
	mysql.DB.Order("time desc").Find(&fileFolders, "parent_folder_id = ? and file_store_id = ?", parentId, fileStoreId)
	return
}

//获取当前的目录信息
func GetCurrentFolders(fId string) (fileFolder data.FileFolder) {
	mysql.DB.Find(&fileFolder, "id = ?", fId)
	return
}

//获取当前路径所有的父级
func GetCurrentAllParent(folder data.FileFolder, folders []data.FileFolder) []data.FileFolder {
	var parentFolder data.FileFolder
	if folder.ParentFolderId != 0 {
		mysql.DB.Find(&parentFolder, "id = ?", folder.ParentFolderId)
		folders = append(folders, parentFolder)
		//递归查找当前所有父级
		return GetCurrentAllParent(parentFolder, folders)
	}

	//反转切片
	for i, j := 0, len(folders)-1; i < j; i, j = i+1, j-1 {
		folders[i], folders[j] = folders[j], folders[i]
	}

	return folders
}

//新建文件夹
func CreateFolder(folderName, parentId string, fileStoreId int) {
	parentIdInt, err := strconv.Atoi(parentId)
	if err != nil {
		fmt.Println("父类id错误")
		return
	}
	fileFolder := data.FileFolder{
		FileFolderName: folderName,
		ParentFolderId: parentIdInt,
		FileStoreId:    fileStoreId,
		Time:           time.Now().Format("2006-01-02 15:04:05"),
	}
	mysql.DB.Create(&fileFolder)
}

//修改文件夹名
func UpdateFolderName(fId, fName string) {
	var fileFolder data.FileFolder
	mysql.DB.Model(&fileFolder).Where("id = ?", fId).Update("file_folder_name", fName)
}

//删除文件夹信息
func DeleteFileFolder(fId string) bool {
	var fileFolder data.FileFolder
	var fileFolder2 data.FileFolder
	//删除文件夹信息
	mysql.DB.Where("id = ?", fId).Delete(data.FileFolder{})
	//删除文件夹中文件信息
	mysql.DB.Where("parent_folder_id = ?", fId).Delete(data.MyFile{})
	//删除文件夹中文件夹信息
	mysql.DB.Find(&fileFolder, "parent_folder_id = ?", fId)
	mysql.DB.Where("parent_folder_id = ?", fId).Delete(data.FileFolder{})

	mysql.DB.Find(&fileFolder2, "parent_folder_id = ?", fileFolder.Id)
	if fileFolder2.Id != 0 { //递归删除文件下的文件夹
		return DeleteFileFolder(strconv.Itoa(fileFolder.Id))
	}

	return true
}

//添加文件数据
func CreateFile(filename, fileHash string, fileSize int64, fId string, fileStoreId int, file_location string) {
	var sizeStr string
	//获取文件后缀
	fileSuffix := path.Ext(filename)
	//获取文件名
	filePrefix := filename[0 : len(filename)-len(fileSuffix)]
	fid, _ := strconv.Atoi(fId)

	if fileSize < 1048576 {
		sizeStr = strconv.FormatInt(fileSize/1024, 10) + "KB"
	} else {
		sizeStr = strconv.FormatInt(fileSize/102400, 10) + "MB"
	}

	myFile := data.MyFile{
		FileName:       filePrefix,
		FileHash:       fileHash,
		FileStoreId:    fileStoreId,
		FilePath:       file_location,
		DownloadNum:    0,
		UploadTime:     time.Now().Format("2006-01-02 15:04:05"),
		ParentFolderId: fid,
		Size:           fileSize / 1024,
		SizeStr:        sizeStr,
		Type:           util.GetFileTypeInt(fileSuffix),
		Postfix:        strings.ToLower(fileSuffix),
	}
	mysql.DB.Create(&myFile)
}
