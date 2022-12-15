package conf

const (

	/*
		accessKey 和 secretKey 是ceph集群搭建后的创建的user生成的，可以通过命令
		docker exec ceph-rgw radosgw-admin user create --uid="user1" --display-name="user1"
		创建用户并查看
	*/
	// CephAccessKey : 访问Key
	CephAccessKey = "0JC8W5QG4MDZ9V8K7NYB"
	// CephSecretKey : 访问密钥
	CephSecretKey = "drcnFT4GVb1NTLAJDTICRZ7ScgAr3icBOejNujRT"
	// CephGWEndpoint : gateway地址
	CephGWEndpoint = "http://127.0.0.1:7000"
)
