package main

import (
	"fmt"
	"io"
	"net"

	pb "github.com/geraldkohn/grpc-example/pb/stream-both"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(stream pb.Greeter_SayHelloServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if ok {
		strs := md.Get("user-id")
		if len(strs) > 0 {
			fmt.Printf("服务器收到元数据: user-id %s\n", strs[0])
		}
	}
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

		fmt.Printf("收到请求: %s\n", in.GetRequst())
		reply := handle(in.GetRequst())

		// 返回流式响应
		if err := stream.Send(&pb.HelloResponse{Reply: reply}); err != nil {
			return err
		}
	}
}

func handle(s string) string {
	// time.Sleep(100 * time.Second) // 模拟处理请求时间
	s = "服务器返回响应: " + s
	return s
}

func main() {
	lis, _ := net.Listen("tcp", ":8972")
	s := grpc.NewServer()                  // 创建grpc服务器
	pb.RegisterGreeterServer(s, &server{}) // 在grpc注册服务
	fmt.Println("Listen: 127.0.0.1:8972")
	_ = s.Serve(lis) // 启动服务
}
