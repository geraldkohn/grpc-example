package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/geraldkohn/grpc-example/pb/stream-client"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8972", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to connect")
		return
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// 执行 grpc 远程调用
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// 客户端流式 RPC
	stream, _ := c.SayHello(ctx)
	reqs := []string{"client-1 ", "client-2 ", "client-3 "}

	// 发送流式数据
	for _, req := range reqs {
		err = stream.Send(&pb.HelloRequest{Requst: req})
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// 接收服务端数据
	res, _ := stream.CloseAndRecv()
	fmt.Println(res)
}
