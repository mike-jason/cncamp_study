# 1.重新修改deployment的yaml，实现了日常运维需求，日志等级
```
k create  -f env_cqhttpserver.yaml

k apply -f deployment_cqhttpserver.yaml
```
# 2.建立ingress控制平面
```
 k create -f nginx-ingress-deployment.yaml
 ```
# 3.生成key和cert
```
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=cqhttp.com/O=cncamp" -addext "subjectAltName = DNS:cqhttp.com"
```
# 4.创建secret对象
```
kubectl create secret tls cqhttp-tls --cert=./tls.crt --key=./tls.key
```
# 5.创建ingress
```
kubectl create -f ingress_cqhttpserver.yaml
```
# 6.service沿用的上个模块中的文件service.yaml
