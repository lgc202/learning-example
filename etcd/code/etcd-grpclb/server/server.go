package main

import (
	pb "code/etcd-grpclb/proto"
	"context"
	"fmt"
	kitLog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/etcdv3"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (g *greeterServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("SayHello Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (g *greeterServer) SayBye(ctx context.Context, in *pb.ByeRequest) (*pb.ByeReply, error) {
	log.Printf("SayBye Received: %v", in.GetName())
	return &pb.ByeReply{Message: "Bye " + in.GetName()}, nil
}

func main() {
	// 动态获取可用的port
	port, err := GetFreePort()
	if err != nil {
		panic(err)
	}

	var (
		etcdServer     = "127.0.0.1:2379"                        // etcd服务的IP地址
		prefix         = "/services/hello/"                      // 服务的目录
		serverInstance = fmt.Sprintf("%s:%d", "127.0.0.1", port) // 当前实例Server的地址
		key            = prefix + serverInstance                 // 服务实例注册的路径
		value          = serverInstance
		ctx            = context.Background()
	)

	// etcd连接参数
	option := etcdv3.ClientOptions{
		DialTimeout:   time.Second * 3,
		DialKeepAlive: time.Second * 3,
	}

	//创建连接
	client, err := etcdv3.NewClient(ctx, []string{etcdServer}, option)
	if err != nil {
		panic(err)
	}

	registrar := etcdv3.NewRegistrar(client, etcdv3.Service{Key: key, Value: value}, kitLog.NewNopLogger())
	registrar.Register() // 启动注册服务

	listen, err := net.Listen("tcp", serverInstance)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &greeterServer{})
	log.Println("Serving gRPC on " + serverInstance)

	go func() {
		if err := s.Serve(listen); err != nil {
			log.Fatalf("listen: %v\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	registrar.Deregister() // 反注册服务
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
