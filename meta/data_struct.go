package meta

import "time"

/*
数据结构和数据库表的映射规则：
假设数据库里存的是 - user_name
那么结构体就应该写 - UserName (因为写成Unername  debug了好久)

如果数据库里存的是 - password
那么结构体就应该是 - password (不区分大小写！)
*/

//用户表
type User struct {
	Id           int
	OpenId       string
	FileStoreId  int
	UserName     string
	PassWord     string
	RegisterTime time.Time
	ImagePath    string
}

//文件存储表
type FileStore struct {
	Id          int
	UserId      int
	CurrentSize int64
	MaxSize     int64
}

//文件表
type MyFile struct {
	Id             int
	FileName       string //文件名
	FileHash       string //文件哈希值
	FileStoreId    int    //文件仓库id
	FilePath       string //文件存储路径
	DownloadNum    int    //下载次数
	UploadTime     string //上传时间
	ParentFolderId int    //父文件夹id
	Size           int64  //文件大小
	SizeStr        string //文件大小单位
	Type           int    //文件类型
	Postfix        string //文件后缀
}

//文件夹表
type FileFolder struct {
	Id             int
	FileFolderName string
	ParentFolderId int
	FileStoreId    int
	Time           string
}

// 分享文件表
type Share struct {
	Id       int
	Code     string
	FileId   int
	Username string
	Hash     string
}
