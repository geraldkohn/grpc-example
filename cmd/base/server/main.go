package main

// base/hello

import (
	"context"
	"fmt"
	"net"

	pb "github.com/geraldkohn/grpc-example/pb/base"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	// 获取元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		strs := md.Get("user-id")
		fmt.Println(strs[0])
	}

	return &pb.HelloResponse{Reply: "hello from server. " + in.GetRequst()}, nil
}

func main() {
	lis, _ := net.Listen("tcp", ":8972")
	s := grpc.NewServer()                  // 创建grpc服务器
	pb.RegisterGreeterServer(s, &server{}) // 在grpc注册服务
	fmt.Println("Listen: 127.0.0.1:8972")
	_ = s.Serve(lis) // 启动服务
}
