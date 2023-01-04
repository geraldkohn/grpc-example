package main

// base/hello

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/geraldkohn/grpc-example/pb/base"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var _debug = true

func debug(format string, a ...interface{}) {
	if _debug {
		s := fmt.Sprintf(format, a...)
		fmt.Println(s)
	}
}

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Request: %v", in)
	log.Printf("Request: %+v", in)
	log.Printf("Request: %#v", in)

	// 获取元数据
	// 将鉴权转移到了拦截器中
	// md, ok := metadata.FromIncomingContext(ctx)
	// if ok {
	// 	strs := md.Get("user-id")
	// 	fmt.Println(strs[0])
	// }
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		user_id := md.Get("user-id")
		debug("user-id: %s", user_id)
	}

	return &pb.HelloResponse{Reply: "hello from server. " + in.GetRequst()}, nil
}

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// token verification
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "没有元数据")
	}
	if !valid(md.Get("user-id")) {
		return nil, status.Errorf(codes.Unauthenticated, "token 失效")
	}

	// 向 handler 传递处理后的结果
	// 新建了 metadata, 并将其传递给 ctx
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("user-id", "经过鉴权的user-id"))
	res, err := handler(ctx, req)
	if err != nil {
		debug("RPC failed with error %v\n", err)
	}
	return res, err
}

func valid(token []string) bool {
	if len(token) <= 0 {
		return false
	}
	debug(token[0])
	return true
}

func main() {
	lis, _ := net.Listen("tcp", ":8972")
	// 创建grpc服务器
	s := grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor),
	)
	pb.RegisterGreeterServer(s, &server{}) // 在grpc注册服务
	fmt.Println("Listen: 127.0.0.1:8972")
	_ = s.Serve(lis) // 启动服务
}
