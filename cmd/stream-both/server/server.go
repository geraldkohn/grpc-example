package main

import (
	"fmt"
	"io"
	"net"

	pb "github.com/geraldkohn/grpc-example/pb/stream-both"
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

func (s *server) SayHello(stream pb.Greeter_SayHelloServer) error {
	// 接收元数据（没使用拦截器）
	// md, ok := metadata.FromIncomingContext(stream.Context())
	// if ok {
	// 	strs := md.Get("user-id")
	// 	if len(strs) > 0 {
	// 		fmt.Printf("服务器收到元数据: user-id %s\n", strs[0])
	// 	}
	// }

	// 处理元数据（使用拦截器后）

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

func streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	debug("进入拦截器")
	// 处理元数据
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return status.Errorf(codes.InvalidArgument, "缺失元数据")
	}
	if !valid(md["user-id"]) {
		return status.Errorf(codes.Unauthenticated, "token 失效")
	}

	// 向 handler 传递处理后的结果
	// 通过 header 来传递信息
	// ss.SendHeader(metadata.Pairs("user-id", "经过鉴权的user-id"))
	
	debug("调用 RPC")
	err := handler(srv, ss)
	debug("调用结束")
	if err != nil {
		debug("RPC failed with error %v\n", err)
	}
	return err
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
	// 创建grpc服务器, 注册拦截器
	s := grpc.NewServer(
		grpc.ChainStreamInterceptor(streamInterceptor), // 链式拦截器
	)
	pb.RegisterGreeterServer(s, &server{}) // 在grpc注册服务
	fmt.Println("Listen: 127.0.0.1:8972")
	_ = s.Serve(lis) // 启动服务
}
