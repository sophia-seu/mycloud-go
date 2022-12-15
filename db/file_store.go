package db

import (
	"mycloud/db/mysql"

	data "mycloud/meta"
)

//根据用户id获取仓库信息
func GetUserFileStore(userId int) (fileStore data.FileStore) {
	mysql.DB.Find(&fileStore, "user_id = ?", userId)
	return
}

//文件上传成功减去相应容量
func SubtractSize(size int64, fileStoreId int) {
	var fileStore data.FileStore
	mysql.DB.First(&fileStore, fileStoreId)

	fileStore.CurrentSize = fileStore.CurrentSize + size/1024
	fileStore.MaxSize = fileStore.MaxSize - size/1024
	mysql.DB.Save(&fileStore)
}

//判断用户容量是否足够
func CapacityIsEnough(fileSize int64, fileStoreId int) bool {
	var fileStore data.FileStore
	mysql.DB.First(&fileStore, fileStoreId)

	return fileStore.MaxSize-(fileSize/1024) >= 0
}
