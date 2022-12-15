package main

// import (
// 	"fmt"
// 	"mycloud/lib"

// 	"gopkg.in/amz.v1/s3"
// )

// func main() {

// 	// bucket := lib.GetCephBucket("admin")

// 	// // d, _ := bucket.Get("/home/sophia/data/ceph/admin/wechat-logo.png")

// 	// // write_location =

// 	// res, _ := bucket.List("", "", "", 100)
// 	// fmt.Println("object keys: ", res)

// 	bucket := lib.GetCephBucket("testbucket1")

// 	// 创建一个新的bucket

// 	err := bucket.PutBucket(s3.PublicRead)
// 	fmt.Println("create bucket err: ", err)

// 	//查询bucket下面制定条件的object keys
// 	res, _ := bucket.List("", "", "", 100)
// 	fmt.Println("object keys", res)
// 	//新上传一个对象
// 	bucket.Put("/home/sophia/data/local/admin/wechat-logo.png", []byte("just for test"), "octet-stream", s3.PublicRead)

// 	// 查询这个bucket下面制定条件的object keys

// 	res, _ = bucket.List("", "", "", 100)
// 	fmt.Println("object keys", res)

// }
