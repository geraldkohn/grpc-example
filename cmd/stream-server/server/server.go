package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/geraldkohn/grpc-example/pb/stream-server"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(in *pb.HelloRequest, stream pb.Greeter_SayHelloServer) error {
	log.Println("Use SayHello Method")
	words := []string{
		"hello",
		"你好",
		"hello again",
		"你好吗",
	}

	for _, word := range words {
		data := &pb.HelloResponse{
			Reply: word + in.GetRequst(),
		}
		// 使用 send 方法返回多个数据
		if err := stream.Send(data); err != nil {
			fmt.Println("服务端错误")
			return err
		}
	}

	return nil
}

func main() {
	lis, _ := net.Listen("tcp", ":8972")
	s := grpc.NewServer()                  // 创建grpc服务器
	pb.RegisterGreeterServer(s, &server{}) // 在grpc注册服务
	fmt.Println("Listen: 127.0.0.1:8972")
	_ = s.Serve(lis) // 启动服务
}
