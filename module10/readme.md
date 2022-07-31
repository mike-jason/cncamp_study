# 1.修改main.go 添加rootHandler函数使用randInt生成随机数0-2000达到延时0-2s的效果
```go
func rootHandler(w http.ResponseWriter, r *http.Request) {
	timer := metrics.NewExecutionTimer()
	defer timer.ObserveTotal()
	user := r.URL.Query().Get("user")
	delay := randInt(0, 2000)
	time.Sleep(time.Millisecond * time.Duration(delay))
	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [#{user}]\n"))
	} else {
		io.WriteString(w, "hello [stranger]\n")
	}
	io.WriteString(w, "====================details of http server header====================\n")
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("#{k}=#{v}\n", k, v))
	}
}
```
# 2.添加metrics目录，编写metrics.go 通过ObserveTotal()函数将延时数据注入到promethues直方图中
```go
//当前调用ObserveTotal时间减去start的时间放到直方图里面
func (e *ExecutionTimer) ObserveTotal() {
	e.histo.WithLabelValues("total").Observe(time.Now().Sub(e.start).Seconds())
}
 ```
# 3.编写deployment.yaml将httpserver部署到k8s测试环境中
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cqhttpserver-metrics
spec:
  replicas: 1
  selector:
    matchLabels:
      app: cqhttpserver-metrics
  template:
    metadata:
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "9999"
      labels:
        app: cqhttpserver-metrics
    spec:
      containers:
        - name: cqhttpserver-metrics
          image: cncamp/cqhttpserver:v1.2-metrics
          ports:
            - containerPort: 9999
```
