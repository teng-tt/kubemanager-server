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

### Pod管理接口开发

- [x] 命名空间列表接口
- Pod创建
- Pod查看（详情、列表）
---
- Pod编辑（更新/升级）
- Pod删除
