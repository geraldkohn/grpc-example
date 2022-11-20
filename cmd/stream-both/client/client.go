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

var _debug = true

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

	// 发送元数据
	// ctx = metadata.AppendToOutgoingContext(ctx, "user-id", "1000000")

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
				debug("客户端发送一条消息")
				req := fmt.Sprintf("Client send at %s", time.Now().String())
				err := stream.Send(&pb.HelloRequest{Requst: req})
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
	debug("客户端停止发送消息")
	sendClose <- struct{}{} // 停止发送消息
	debug("客户端关闭连接")
	stream.CloseSend() // 客户端关闭连接
	debug("等待客户端收到服务器全部消息")
	<-wait // 等待客户端接收到全部信息
	debug("收到服务器全部消息，结束")
}

func debug(format string, a ...interface{}) {
	if _debug {
		s := fmt.Sprintf(format, a...)
		fmt.Println(s)
	}
}
