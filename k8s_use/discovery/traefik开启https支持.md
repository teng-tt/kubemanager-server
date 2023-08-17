

# traefik开启https支持


## 修改Traefik服务配置

args添加参数
```bash
- --entrypoints.http.address=:80
- --entrypoints.https.address=:443
```

ports新增端口映射：

```bash
- name: https
  containerPort: 443
  hostPort: 443
```

## 生成证书

```bash
openssl req -newkey rsa:2048 -nodes -keyout tls.key -x509 -days 3650 -out tls.crt
```

## 创建Secret

```bash
kubectl create secret generic kubemanager-cert --from-file=tls.crt --from-file=tls.key -n kubemanager-system
```

## IngressRoute指定secret
```yaml
tls:
  secretName: kubemanager-cert
```