package main

import (
	pb "code/etcd-grpclb/proto"
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	kitLog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
	"google.golang.org/grpc"
	"io"
	"strconv"
	"time"
)

const (
	DefaultName = "casa"
	SayHello    = "SayHello"
	SayByte     = "SayByte"
)

type EndPointRequest struct {
	Method string
	Req    interface{}
}

func main() {

	var (
		etcdServer = "127.0.0.1:2379"   // 注册中心地址
		prefix     = "/services/hello/" // 监听的服务前缀
		ctx        = context.Background()
	)

	// etcd连接参数
	options := etcdv3.ClientOptions{
		DialTimeout:   time.Second * 3,
		DialKeepAlive: time.Second * 3,
	}

	// 连接注册中心
	client, err := etcdv3.NewClient(ctx, []string{etcdServer}, options)
	if err != nil {
		panic(err)
	}

	// 创建实例管理器, 此管理器会Watch监听etc中prefix的目录变化更新缓存的服务实例数据
	logger := kitLog.NewNopLogger()
	instancer, err := etcdv3.NewInstancer(client, prefix, logger)
	if err != nil {
		panic(err)
	}

	// 创建端点管理器， 此管理器根据Factory和监听的到实例创建endPoint并订阅instancer的变化动态更新Factory创建的endPoint
	endpointer := sd.NewEndpointer(instancer, reqFactory, logger) // reqFactory 是自定义的业务处理函数

	// 创建轮询负载均衡器
	balancer := lb.NewRoundRobin(endpointer)

	// 获取Endpoint用来进行grpc调用
	retry := lb.Retry(3, 3*time.Second, balancer)

	// 模拟调用10次SayHello和10次SayByte，这20个请求依次轮询发送到不同的服务端
	for i := 0; i < 10; i++ {
		// 调用SayHello方法
		helloRequest := EndPointRequest{
			Method: SayHello,
			Req:    &pb.HelloRequest{Name: DefaultName + " " + strconv.Itoa(i)},
		}
		rsp, err := retry(ctx, helloRequest)
		if err != nil {
			fmt.Println(err)
			continue
		}
		helloReply, _ := rsp.(*pb.HelloReply)
		fmt.Println("Greeting: ", helloReply.GetMessage())

		// 调用SayByte方法
		byeRequest := EndPointRequest{
			Method: SayByte,
			Req:    &pb.ByeRequest{Name: DefaultName + " " + strconv.Itoa(i)},
		}
		byeRsp, err := retry(ctx, byeRequest)
		if err != nil {
			fmt.Println(err)
			continue
		}
		byeReply, _ := byeRsp.(*pb.ByeReply)
		fmt.Println("Greeting: ", byeReply.GetMessage())

		time.Sleep(time.Second)
	}
}

// 通过传入的实例地址  创建对应的请求endPoint
func reqFactory(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("请求服务: ", instanceAddr)
		conn, err := grpc.Dial(instanceAddr, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		client := pb.NewGreeterClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		req, _ := request.(EndPointRequest)
		switch req.Method {
		case SayHello:
			helloReq, _ := req.Req.(*pb.HelloRequest)
			return client.SayHello(ctx, helloReq)
		case SayByte:
			byeReq, _ := req.Req.(*pb.ByeRequest)
			return client.SayBye(ctx, byeReq)
		default:
			return nil, fmt.Errorf("unsupport method: %s", req.Method)
		}
	}, nil, nil
}
