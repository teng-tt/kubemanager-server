# kubemanager-server

go语言开发的的k8s集群管理工具，基于web可视化的方式对k8s集群进行管理，简化k8s集群操作难度 \
同时整合Harbor与 prometheus实现镜像管理和仪表盘监控

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

### 7. k8s认证与授权
在集群内初始化，不需要知道.kube/config

ServiceAccount
- [x] 创建
- [x] 删除
- [x] 查询(列表)

Role/ClusterRole
- [x] 创建/更新
- [x] 删除
- [x] 查询(详情/列表)

RoleBinding/ClusterRoleBinding
- [x] 创建/更新
- [x] 删除
- [x] 查询(详情/列表)

### 8. k8s数据备份crd+operator开发
通过kubebuilder脚手架开发数据备份服务 \
实现按指定时间与间隔，将数据源备份到目标存储源 \
此处配置的数据源为mysql8.0, 目标存储源为oss 

### 9. kubemanager-server 整合 Harbor
在管理平台整合Harbor v2实现容器镜像管理
- [x] 集成HarborAPI
- [x] Projects 列表查询(分页|模糊查询)
- [x] Repositories 列表查询(分页|模糊查询)
- [x] Artifacts 列表查询(分页|模糊查询)
- [x] 镜像匹配接口(用户Pod输入镜像信息时，自动匹配)

### 10. 仪表盘功能
在管理平台实现监控信息展示，整合prometheus
- [x] 基础信息查看(k8s版本信息，集群初始化时间等)
- [x] 各资源的统计信息
- [x] 集群pod、cpu、内存耗用情况(瞬时)
    - [x] 安装metrics-server
    - [x] 调用metrics-server接口，计算集群的cpu和内存的耗用
- [x] 集群 cpu、内存变化趋势
    - [x] 安装prometheus
    - [x] 提供prometheus pull 数据的接口(exporter)
    - [x] 调用提供prometheus 查询指标统计数据