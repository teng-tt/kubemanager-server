# nfs环境搭建

简介 NFS（Network File System）**即网络文件系统，是FreeBSD支持的文件系统中的一种，它允许网络中的计算机之间共享资源**。 在NFS的应用中，本地NFS的客户端应用可以透明地读写位于远端NFS服务器上的文件，就像访问本地文件一样。

## 安装nfs（nfs主节点）

```bash
sudo apt-get install nfs-kernel-server  # 安装 NFS服务器端
sudo apt-get install nfs-common         # 安装 NFS客户端
```

> nfs服务器端，只在集群中某一台安装即可

## 配置

创建nfs共享目录

```bash
mkdir -p /nfs/share
chmod -R 666 /nfs/share
```

修改配置 `vim /etc/exports`

```bash
#若需要把 “/data/backups” 目录设置为 NFS 共享目录，请在该文件末尾添加下面的一行：
# 当登录NFS主机使用共享目录的使用者是root时，其权限将被转换成为匿名使用者，通常它的UID与GID都会变成nobody身份,添加no_root_squash参数，确保root账户能用
/nfs/share *(rw,sync,no_root_squash)     # * 表示允许任何网段 IP 的系统访问该 NFS 目录
```

配置生效
```
exportfs -r
```

查看生效

```bash
exportfs
```

## 启动服务

```bash
sudo /etc/init.d/nfs-kernel-server start    或者  
sudo /etc/init.d/nfs-kernel-server restart
```

## 测试挂载（nfs node节点）

```bash
mkdir -p /nfs/share
mount -t nfs 192.168.1.16:/nfs/share /nfs/share -o nolock
```

在挂载的目录写入数据，则能看到主节点也能同步写入。

如果要解绑，则执行命令：

```bash
umount /nfs/share
```