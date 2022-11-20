package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "github.com/geraldkohn/grpc-example/pb/base"
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
	// 添加元数据
	ctx = metadata.AppendToOutgoingContext(ctx, "user-id", "10000000")
	r, _ := c.SayHello(ctx, &pb.HelloRequest{Requst: "client"})
	fmt.Println(r.GetReply())
}
