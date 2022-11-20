package main

import (
	"fmt"
	"io"
	"net"
	"time"

	pb "github.com/geraldkohn/grpc-example/pb/stream-both"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(stream pb.Greeter_SayHelloServer) error {
	for {
		// 接收流式请求
		in, err := stream.Recv()
		// 都客户端关闭发送连接, 就退出
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		reply := handle(in.GetRequst())

		// 返回流式响应
		if err := stream.Send(&pb.HelloResponse{Reply: reply}); err != nil {
			return err
		}
	}
}

func handle(s string) string {
	time.Sleep(100 * time.Second) // 模拟处理请求时间
	s += " Server Response"
	return s
}

func main() {
	lis, _ := net.Listen("tcp", ":8972")
	s := grpc.NewServer()                  // 创建grpc服务器
	pb.RegisterGreeterServer(s, &server{}) // 在grpc注册服务
	fmt.Println("Listen: 127.0.0.1:8972")
	_ = s.Serve(lis) // 启动服务
}
