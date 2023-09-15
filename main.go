package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"grpc_go/codes"

	"google.golang.org/grpc"
)

var grpc_port string
var http_port string
var log_file string
var syncmap sync.Map

type server struct{
    brpc_metrics.UnimplementedMetricServiceServer
}

func (s *server) CollectMetrics(ctx context.Context, req *brpc_metrics.MetricRequest) (*brpc_metrics.MetricResponse, error) {
	for _, item := range req.Metric {
		syncmap.Store(item.Key, item.Value)
	}
	return &brpc_metrics.MetricResponse{Msg: "success"}, nil
}

func init() {
	flag.StringVar(&grpc_port, "grpc_port", ":50033", "port, examle: :7777")
	flag.StringVar(&http_port, "http_port", ":8013", "http port, examle: :777")
	flag.StringVar(&log_file, "logfile", "/data/log/party_brpc_metrics/party_brpc_metrics.DEBUG", "brpc log file")
}

func start_grpc() {
	lis, err := net.Listen("tcp", grpc_port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	brpc_metrics.RegisterMetricServiceServer(s, &server{}) // 注意这里
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	flag.Parse()

	file, err := os.OpenFile(log_file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
    if err != nil {
        log.Fatal("无法打开日志文件: ", err)
    }
	log.SetOutput(file)

    defer file.Close()

	log.Println(grpc_port)
	log.Println(http_port)
	go start_grpc()
	// 创建一个处理器函数，用于处理HTTP请求
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		// 设置响应的内容类型为文本/plain
		w.Header().Set("Content-Type", "text/plain")
		// 写入响应消息

		var metrics string
		// 遍历 sync.Map
		syncmap.Range(func(key, value interface{}) bool {
			str_key := fmt.Sprintf("%v", key)
			str_value := fmt.Sprintf("%v", value)
			metrics += str_key + " " + str_value + "\n"
			return true // 返回true继续遍历，返回false停止遍历
		})
		fmt.Fprintf(w, "%s", metrics)
	})

	err = http.ListenAndServe(http_port, nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

