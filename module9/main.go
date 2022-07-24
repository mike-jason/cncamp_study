package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	http.HandleFunc("/readlog", readlogHandler)
	http.HandleFunc("/readOSvariables", readOSvariablesHandler)
	http.HandleFunc("/responheader", responheaderHandler)
	http.HandleFunc("/healthz", healthz)
	log.Fatal(http.ListenAndServe(":9999", nil))
}

//读取客户端ip (http返回码不知道取哪个字段值，烦请老师指教)
func readlogHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("Client Address:%q\n", r.RemoteAddr)

}

//读取系统环境变量并写入header
func readOSvariablesHandler(w http.ResponseWriter, r *http.Request) {
	VERSION := os.Getenv("VERSION")
	w.Header().Add("version", VERSION)
}

//接收客户端 request，并将 request 中带的 header 写入 response header
func responheaderHandler(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		w.Header().Add(k, strings.Join(v, ""))
	}
}

//当访问 localhost/healthz 时，应返回 200
func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v\n", http.StatusOK)
}

// 优雅终止
func listenSignal(ctx context.Context, httpSrv *http.Server) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-sigs:
		fmt.Println("notify sigs")
		httpSrv.Shutdown(ctx)
		fmt.Println("http shutdown gracefully")
	}
}
