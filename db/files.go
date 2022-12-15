package db

import (
	"mycloud/db/mysql"
	data "mycloud/meta"
	"path"
	"strings"
)

// 获取用户文件
func GetUserFiles(parentId string, storeId int) (files []data.MyFile) {
	mysql.DB.Find(&files, "file_store_id = ? and parent_folder_id = ?", storeId, parentId)
	return
}

//获取用户文件数量
func GetUserFileCount(fileStoreId int) (fileCount int) {
	var file []data.MyFile
	mysql.DB.Find(&file, "file_store_id = ?", fileStoreId).Count(&fileCount)
	return
}

//获取用户文件使用明细情况
func GetFileDetailUse(fileStoreId int) map[string]int64 {
	var files []data.MyFile
	var (
		docCount   int64
		imgCount   int64
		videoCount int64
		musicCount int64
		otherCount int64
	)

	fileDetailUseMap := make(map[string]int64, 0)

	//文档类型
	docCount = mysql.DB.Find(&files, "file_store_id = ? AND type = ?", fileStoreId, 1).RowsAffected
	fileDetailUseMap["docCount"] = docCount
	////图片类型
	imgCount = mysql.DB.Find(&files, "file_store_id = ? and type = ?", fileStoreId, 2).RowsAffected
	fileDetailUseMap["imgCount"] = imgCount
	//视频类型
	videoCount = mysql.DB.Find(&files, "file_store_id = ? and type = ?", fileStoreId, 3).RowsAffected
	fileDetailUseMap["videoCount"] = videoCount
	//音乐类型
	musicCount = mysql.DB.Find(&files, "file_store_id = ? and type = ?", fileStoreId, 4).RowsAffected
	fileDetailUseMap["musicCount"] = musicCount
	//其他类型
	otherCount = mysql.DB.Find(&files, "file_store_id = ? and type = ?", fileStoreId, 5).RowsAffected
	fileDetailUseMap["otherCount"] = otherCount

	return fileDetailUseMap
}

//根据文件类型获取文件
func GetTypeFile(fileType, fileStoreId int) (files []data.MyFile) {
	mysql.DB.Find(&files, "file_store_id = ? and type = ?", fileStoreId, fileType)
	return
}

//通过fileId获取文件信息
func GetFileInfo(fId string) (file data.MyFile) {
	mysql.DB.First(&file, fId)
	return
}

//删除数据库文件数据
func DeleteUserFile(fId, folderId string, storeId int) {
	mysql.DB.Where("id = ? and file_store_id = ? and parent_folder_id = ?",
		fId, storeId, folderId).Delete(data.MyFile{})
}

//文件下载次数+1
func DownloadNumAdd(fId string) {
	var file data.MyFile
	mysql.DB.First(&file, fId)
	file.DownloadNum = file.DownloadNum + 1
	mysql.DB.Save(&file)
}

//判断当前文件夹是否有同名文件
func CurrFileExists(fId, filename string) bool {
	var file data.MyFile
	//获取文件后缀
	fileSuffix := strings.ToLower(path.Ext(filename))
	//获取文件名
	filePrefix := filename[0 : len(filename)-len(fileSuffix)]

	mysql.DB.Find(&file, "parent_folder_id = ? and file_name = ? and postfix = ?", fId, filePrefix, fileSuffix)

	return file.Size <= 0
}
