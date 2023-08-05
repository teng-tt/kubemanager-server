# kubmanager

go语言开发的的k8s集群管理工具

开发环境说明： 
- go语言版本：go1.19 windows/amd64
- 编译环境：Linux
- 开发工具：goland

## 项目的初始化

### web框架的选型
```bash
    go get -u github.com/gin-gonic/gin@v1.8.1
```

### 配置参数分离
```bash
    go get github.com/spf13/viper@v1.13.0
```

> 参考文档：https://github.com/spf13/viper

### k8s工具集成
```bash
    go get k8s.io/client-go@v0.20.4
```

> 参考文档：https://github.com/kubernetes/client-go

## 项目接口开发

### 1. Pod管理接口

- [x] 命名空间列表接口
- [x] Pod创建
- [x] Pod编辑 （更新/升级）
---
- [x] Pod查看-详情
    展示 podrequest 数据 用于重新创建
- [x] Pod查看-列表
- [x] Pod删除
---
接口调优：
- pod更新会多出一个挂载卷 （serviceAccount） 
  - 计算那些是emtydir volume mount 进行非emtydir过滤
- 更新pod超时
  -  pod删除等待时间不确定，改为强制删除，减少删除等待时间，防止前端删除超时
- pod列表支持关键字搜索

### 2. NodeScheduling接口
- [x] node列表/详情(kubectl get nodes / kubectl describe node -node-x)
- [x] node标签管理(kubectl label node node-x label-x=label-value-x)
    - 所有的标签上传
- [x] node污点(taint)管理
- [x] 查看node上所有的Pod(kubectl get pod -n  ns-x -o wide)
---
pod管理接口改动:
- [x] pod新增容忍(tolerations)参数
- [x] pod选择哪种调度方式：nodeName/nodeSelector/nodeAffinity

### 3. 应用与配置分离接口
ConfigMap 
- [x] 新增|修改
- [x] 查询（列表|详情）
- [x] 删除

Secret 
- [x] 新增|修改
- [x] 查询(列表|详情)
- [x] 删除
---
Pod管理接口改动：
- [x] 新增ConfigMap和ConfigMapKeY
- [x] 新增Secret和SecretKey

### 4. k8s卷管理接口
PersistentVolume
- [x] 创建
- [x] 删除
- [x] 查询--列表 

PersistentVolumeClaim
- [x] 创建
- [x] 删除
- [x] 查询--列表 

StorageClass
- [x] 创建
- [x] 删除
- [x] 查询--列表 

pod管理接口改动:
- [x] Pod管理 (卷管理部分逻辑新增存储卷支持，支持`emptydir|confiMap|secret|hostPath|downward|pvc`)
---
优化点：
- [x] downward fileRefPath没有显示
- [x] PVC选择PV或SC只能二选一
- [x] SC PVC PV 添加keyword搜索字段
- [x] PV显示StorageClassName字段

### 5. k8s服务发现接口

Service
- [x] 创建/更新
- [x] 删除
- [x] 查询 (列表和详情)

Ingress
- [x] 创建/更新
- [x] 删除
- [x] 查询 (列表和详情)

IngressRoute(traefik的自定义资源)
- [x] 创建/更新
- [x] 删除
- [x] 查询 (列表和详情)
- [x] Middleware的查询接口
> 自定义资源接口详情参考文档： https://github.com/kubernetes-client/python/blob/master/kubernetes/README.md

### 6. k8s工作负载

StatefulSet
- [x] 创建/更新
- [x] 删除
- [x] 查询 (列表和详情)

Deployment
- [x] 创建/更新
- [x] 删除
- [x] 查询 (列表和详情)

DaemonSet
- [x] 创建/更新
- [x] 删除
- [x] 查询 (列表和详情)

Job
- [x] 创建/更新
- [x] 删除
- [x] 查询 (列表和详情)

CronJob
- [x] 创建/更新
- [x] 删除
- [x] 查询 (列表和详情)