FROM golang:1.17 AS build
WORKDIR /cqhttpserver/
COPY . .
ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
RUN GOOS=linux go build -installsuffix cgo -o cqhttpserver main.go

FROM busybox
COPY --from=build /cqhttpserver/cqhttpserver /cqhttpserver/cqhttpserver
EXPOSE 9999
ENV ENV local
WORKDIR /cqhttpserver/
ENTRYPOINT ["./cqhttpserver"]