package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/geraldkohn/grpc-example/pb/stream-server"
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
	stream, _ := c.SayHello(ctx, &pb.HelloRequest{Requst: "client"})

	for {
		// 接收服务端返回的流式数据，当收到io.EOF或者错误时退出
		res, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("传输结束")
			break
		}
		if err != nil {
			fmt.Println(err)
			fmt.Println("服务端错误")
			break
		}
		fmt.Println(res)
	}
}
