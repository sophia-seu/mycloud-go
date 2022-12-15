通过docker可以快速部署小规模Ceph集群的流程，可用于开发测试。
以下的安装流程是通过linux shell来执行的；假设你只有一台机器，装了linux(如Ubuntu)系统和docker环境，那么可以参考以下步骤安装Ceph:

```bash
# 要用root用户创建, 或有sudo权限
# 注: 建议使用这个docker镜像源:https://registry.docker-cn.com
# 1. 修改docker镜像源
cat > /etc/docker/daemon.json << EOF
{
  "registry-mirrors": [
    "https://registry.docker-cn.com"
  ]
}
EOF
# 重启docker
systemctl restart docker
# 2. 创建Ceph专用网络
docker network create --driver bridge --subnet 172.20.0.0/16 ceph-network
docker network inspect ceph-network
# 3. 清理旧的ceph相关目录文件，加入有的话
rm -rf /www/ceph /var/lib/ceph/  /www/osd/


# 4. 拉取ceph镜像
docker pull ceph/daemon:latest-luminous
# 5. 创建monitor节点
docker run -d --name ceph-mon --network ceph-network --ip 172.20.0.10 -e CLUSTER=ceph -e WEIGHT=1.0 -e MON_IP=172.20.0.10 -e MON_NAME=ceph-mon -e CEPH_PUBLIC_NETWORK=172.20.0.0/16 -v /etc/ceph:/etc/ceph -v /var/lib/ceph/:/var/lib/ceph/ -v /var/log/ceph/:/var/log/ceph/ ceph/daemon:latest-luminous mon

# 6. 搭建osd节点
docker exec ceph-mon ceph auth get client.bootstrap-osd -o /var/lib/ceph/bootstrap-osd/ceph.keyring

# 7. 修改配置文件以兼容etx4硬盘
vi /etc/ceph/ceph.conf 
# 在文件最后追加
osd max object name len = 256
osd max object namespace len = 64

# 8. 分别启动三个容器来模拟集群
docker run -d --privileged=true --name ceph-osd-1 --network ceph-network --ip 172.20.0.11 -e CLUSTER=ceph -e WEIGHT=1.0 -e MON_NAME=ceph-mon -e MON_IP=172.20.0.10 -e OSD_TYPE=directory -v /etc/ceph:/etc/ceph -v /var/lib/ceph/:/var/lib/ceph/ -v /var/lib/ceph/osd/1:/var/lib/ceph/osd -v /etc/localtime:/etc/localtime:ro ceph/daemon:latest-luminous osd

docker run -d --privileged=true --name ceph-osd-2 --network ceph-network --ip 172.20.0.12 -e CLUSTER=ceph -e WEIGHT=1.0 -e MON_NAME=ceph-mon -e MON_IP=172.20.0.10 -e OSD_TYPE=directory -v /etc/ceph:/etc/ceph -v /var/lib/ceph/:/var/lib/ceph/ -v /var/lib/ceph/osd/2:/var/lib/ceph/osd -v /etc/localtime:/etc/localtime:ro ceph/daemon:latest-luminous osd

docker run -d --privileged=true --name ceph-osd-3 --network ceph-network --ip 172.20.0.13 -e CLUSTER=ceph -e WEIGHT=1.0 -e MON_NAME=ceph-mon -e MON_IP=172.20.0.10 -e OSD_TYPE=directory -v /etc/ceph:/etc/ceph -v /var/lib/ceph/:/var/lib/ceph/ -v /var/lib/ceph/osd/3:/var/lib/ceph/osd -v /etc/localtime:/etc/localtime:ro ceph/daemon:latest-luminous osd

# 9. 创建gateway节点
docker exec ceph-mon ceph auth get client.bootstrap-rgw -o /var/lib/ceph/bootstrap-rgw/ceph.keyring

docker run -d --privileged=true --name ceph-mgr --network ceph-network --ip 172.20.0.14 -e CLUSTER=ceph -p 7000:7000 --pid=container:ceph-mon -v /etc/ceph:/etc/ceph -v /var/lib/ceph/:/var/lib/ceph/ ceph/daemon:latest-luminous mgr
# 10. 查看ceph集群状态
sleep 10 && docker exec monnode ceph -s
# 11. 创建用户
docker exec -it gwnode radosgw-admin user create --uid=user1 --display-name=user1
```
