package main

import (
	"fmt"
	"io"
	"net"

	pb "github.com/geraldkohn/grpc-example/pb/stream-client"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(stream pb.Greeter_SayHelloServer) error {
	reply := "hello, this is server, request from client: "
	res := ""
	for {
		// 接收客户端发来的流数据
		req, err := stream.Recv()
		res += req.GetRequst()
		if err == io.EOF {
			return stream.SendAndClose(&pb.HelloResponse{
				Reply: fmt.Sprintf("%s %s", reply, res),
			})
		}
		if err != nil {
			return err
		}
	}
}

func main() {
	lis, _ := net.Listen("tcp", ":8972")
	s := grpc.NewServer()                  // 创建grpc服务器
	pb.RegisterGreeterServer(s, &server{}) // 在grpc注册服务
	fmt.Println("Listen: 127.0.0.1:8972")
	_ = s.Serve(lis) // 启动服务
}
