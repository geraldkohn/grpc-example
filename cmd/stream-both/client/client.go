package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/geraldkohn/grpc-example/pb/stream-both"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 双向流式 RPC
	stream, _ := c.SayHello(ctx)

	wait := make(chan struct{})
	sendClose := make(chan struct{})

	// 发送流数据
	go func() {
		for {
			select {
			case <-sendClose:
				return
			default:
				err := stream.Send(&pb.HelloRequest{Requst: fmt.Sprintf("Client send at %s", time.Now().String())})
				if err != nil {
					return
				}
				time.Sleep(10 * time.Millisecond) // 每隔 10 毫秒发送一个信息
			}
		}
	}()

	// 接收流数据
	go func() {
		for {
			res, err := stream.Recv()
			// 已经接收到全部消息
			if err == io.EOF {
				wait <- struct{}{} // 服务器处理完毕
				return
			}
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(res.GetReply())
		}
	}()

	time.Sleep(1 * time.Second)
	sendClose <- struct{}{} // 停止发送消息
	stream.CloseSend()      // 客户端关闭连接
	<-wait                  // 等待客户端接收到全部信息
}
