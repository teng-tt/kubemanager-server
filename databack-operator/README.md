# databack-operator
数据备份，基于k8s crd+operator实现

环境说明：
- 开发环境： windows11 go1.19
- 运行环境： linux amd64(k8s master--kubectl环境，操作k8s)

## 项目初始化

基于kubebuilder完成初始化(kubebuilder是一个crd+operator脚手架)
```bash 
https://github.com/kubernetes-sigs/kubebuilder.git
```
编译kubebuilder
```bash
make build
```

基于kubebuilder初始化代码结构
```bash
kubebuilder init databack-operator --domain="operator.kubemanager.com" --project-name="databack-operator" --repo="kubemanager.com/operator-databackup"
```

创建api
```bash
kubebuilder create api --group "" --version v1beta1 --kind Databack
```

## 项目发布

安装crd
```bash
make install
```

打包operator
```bash
make docker-build docker-push IMG=xx.com/operator/databack:v1beta1
```

发布operator
```bash
make deploy IMG=xx.com/operator/databack:v1beta1
```

部署databack服务
```bash
kubectl apply -f config/samples_v1beta1_databack.yaml
```

## 项目开发

### 定义crd
修改： /api/v1beta1/databack_types.go

### 开发operator
修改: /controllers/databack_controller.go
实现Reconcile的逻辑

## 正式对外提供服务
通过下面得到的yaml文件，就可以去所有的K8S环境运行部署实例提供服务
```bash
/databack-operator/bin/kustomize build config/default > databack-operator.yaml
```
