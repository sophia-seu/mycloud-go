package db

import (
	"log"
	"mycloud/db/mysql"
	data "mycloud/meta"
	"time"
)

// 创建用户
func CreateUser(openId, username, password, imageUrl string) {
	// 先让FileStoreId为0，然后创建在file_store表中增加一行
	// 再将file_store中的Id值赋回FileStoreId，由于Id是自增的，所以可以自增FileStoreId
	user := data.User{
		OpenId:       openId,
		FileStoreId:  0,
		UserName:     username,
		PassWord:     password,
		RegisterTime: time.Now(),
		ImagePath:    imageUrl,
	}

	if err := mysql.DB.Create(&user).Error; err != nil {
		log.Println("Falied to insert new user...")
		return
	}

	fileStore := data.FileStore{
		UserId:      user.Id,
		CurrentSize: 0,
		MaxSize:     1048576,
	}
	if err := mysql.DB.Create(&fileStore).Error; err != nil {
		log.Println("Falied to insert to filestore...")
		return
	}

	user.FileStoreId = fileStore.Id

	mysql.DB.Save(&user)
}

//查询判断用户是否存在,如果存在，同时返回openID
func QueryUserExists(username, password string) string {
	var user data.User
	mysql.DB.Find(&user, "user_name = ?", username)

	//用户不存在
	if user.Id == 0 {
		return "not exist"
	}

	if user.PassWord != password {
		return "not match"
	}

	return user.OpenId
}

//根据openId查询用户
func GetUserInfo(openId interface{}) (user data.User) {
	mysql.DB.Find(&user, "open_id = ?", openId)
	return
}
