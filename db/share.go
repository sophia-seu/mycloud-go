package db

import (
	"mycloud/db/mysql"
	data "mycloud/meta"
	"mycloud/util"
	"strings"
	"time"
)

//创建分享
func CreateShare(code, username string, fId int) string {
	share := data.Share{
		Code:     strings.ToLower(code),
		FileId:   fId,
		Username: username,
		Hash:     util.EncodeMd5(code + string(time.Now().Unix())),
	}
	mysql.DB.Create(&share)

	return share.Hash
}

//查询分享
func GetShareInfo(f string) (share data.Share) {
	mysql.DB.Find(&share, "hash = ?", f)
	return
}

//校验提取码
func VerifyShareCode(fId, code string) bool {
	var share data.Share
	mysql.DB.Find(&share, "file_id = ? and code = ?", fId, code)
	return share.Id != 0
}
