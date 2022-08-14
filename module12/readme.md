# 作业
把我们的 httpserver 服务以 Istio Ingress Gateway 的形式发布出来。以下是你需要考虑的几点：

- 如何实现安全保证；
- 七层路由规则；
- 考虑 open tracing 的接入。


#创建对应ns并打上istio label
```shell
kubectl create ns jiajia; kubectl label ns jiajia istio-injection=enabled
```
#生成key和crt
```shell
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=cncamp Inc./CN=*.cncamp.io' -keyout cncamp.io.key -out cncamp.io.crt
```
#创建对应的secret(ns: istio-system)
```shell
kubectl create -n istio-system secret tls cncamp-credential --key=cncamp.io.key --cert=cncamp.io.crt
```
#创建服务svc
```shell
kubectl create -f cqhttpserver.yaml -n jiajia
```
#创建istio vs和gateway
```shell
k create -f istio-spec.yaml -n jiajia
```
#验证
```shell
curl --resolve cqhttpsserver.cncamp.io:443:10.107.168.179 https://cqhttpsserver.cncamp.io/healthz -k
root@vm-0-9-ubuntu:/home/ubuntu/homework/module12# curl --resolve cqhttpsserver.cncamp.io:443:10.107.168.179 https://cqhttpsserver.cncamp.io/healthz -k
200
```
#7层转发规则
```shell
k create -f istio-l7.yaml -n jiajia
```
#验证
```shell
curl --resolve cqhttpsserver.cncamp.io:443:10.107.168.179 https://cqhttpsserver.cncamp.io/l7/hello -k
root@vm-0-9-ubuntu:/home/ubuntu/homework/module12# curl --resolve cqhttpsserver.cncamp.io:443:10.107.168.179 https://cqhttpsserver.cncamp.io/l7/hello -k
hello [stranger]
====================details of http server header====================
#{k}=#{v}
%!(EXTRA string=Accept, []string=[*/*])#{k}=#{v}
%!(EXTRA string=X-Forwarded-Proto, []string=[https])#{k}=#{v}
%!(EXTRA string=X-Envoy-Attempt-Count, []string=[1])#{k}=#{v}
%!(EXTRA string=X-Envoy-Internal, []string=[true])#{k}=#{v}
%!(EXTRA string=X-B3-Spanid, []string=[043215814aeca80b])#{k}=#{v}
%!(EXTRA string=X-B3-Traceid, []string=[911568c846338ea46e10c7c3480f3010])#{k}=#{v}
%!(EXTRA string=X-B3-Parentspanid, []string=[6e10c7c3480f3010])#{k}=#{v}
%!(EXTRA string=X-B3-Sampled, []string=[1])#{k}=#{v}
%!(EXTRA string=User-Agent, []string=[curl/7.68.0])#{k}=#{v}
%!(EXTRA string=X-Forwarded-For, []string=[172.21.0.9])#{k}=#{v}
%!(EXTRA string=X-Request-Id, []string=[58002232-43d2-97c3-a818-7d8f34638ac3])#{k}=#{v}
%!(EXTRA string=X-Envoy-Original-Path, []string=[/l7/hello])#{k}=#{v}
%!(EXTRA string=X-Forwarded-Client-Cert, []string=[By=spiffe://cluster.local/ns/jiajia/sa/default;Hash=f9fb22c2eb14554c4385624a4c0d8fce48443596af5c8a68042d2bb2ff9585b1;Subject="";URI=spiffe://cluster.local/ns/istio-system/sa/istio-ingressgateway-service-account])
```
#Tracing
```shell
kubectl apply -f jaeger.yaml
```
