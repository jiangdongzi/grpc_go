package main

import (
    "context"
    "fmt"
    "log"

    "google.golang.org/grpc"

    pb "grpc_go/codes" // 替换为你的实际 proto 文件所在的包路径
)

func main() {
    // 连接到 gRPC 服务器
    conn, err := grpc.Dial("localhost:50033", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("无法连接到服务器: %v", err)
    }
    defer conn.Close()

    // 创建 MetricService 的客户端
    client := pb.NewMetricServiceClient(conn)

    // 创建一个 MetricRequest
    request := &pb.MetricRequest{
        Metric: []*pb.Metric{
            {
                Key:   "key1",
                Value: "value1",
            },
            {
                Key:   "key2",
                Value: "value2",
            },
        },
    }

    // 调用 CollectMetrics RPC
    response, err := client.CollectMetrics(context.Background(), request)
    if err != nil {
        log.Fatalf("RPC 调用失败: %v", err)
    }

    // 处理 MetricResponse
    fmt.Printf("收到服务器响应: %s\n", response.Msg)
}

