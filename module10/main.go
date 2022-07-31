package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"metrics"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	metrics.RegisterMetrics()
	http.HandleFunc("/readlog", readlogHandler)
	http.HandleFunc("/hello", rootHandler)
	http.HandleFunc("/readOSvariables", readOSvariablesHandler)
	http.HandleFunc("/responheader", responheaderHandler)
	http.HandleFunc("/healthz", healthz)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9999", nil))
}
func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

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

//读取客户端ip
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
