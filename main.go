package main

import (
	RemoteProxy "amqp-proxy/proto"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
)

type server struct {
}

func (s *server) mustEmbedUnimplementedConsumerServer() {
	panic("not implemented")
}

var phpPath string
var artisanCommand = "kyy:proxy"
var logger int
var ipv4 string
var port string

func logConsole(format string, v ...any) {
	if logger == 1 {
		log.Printf(format, v...)
	}
}

func (s *server) Work(ctx context.Context, params *RemoteProxy.Params) (*RemoteProxy.Reply, error) {
	logConsole("%s %s/artisan %s %s %s %s", phpPath, params.GetPath(), artisanCommand, params.GetClass(), params.GetFunc(), params.GetArgs())
	go func() {
		cmd := exec.Command(phpPath, params.GetPath()+"/artisan", artisanCommand, params.GetClass(), params.GetFunc(), "--args="+params.GetArgs())
		out, _ := cmd.CombinedOutput()
		logConsole("combined out:\n%s\n", string(out))
	}()
	return &RemoteProxy.Reply{Success: true, Message: ""}, nil
}

func getClientIp() string {
	host := "0.0.0.0"
	adds, err := net.InterfaceAddrs()
	if err != nil {
		return host
	}
	for _, address := range adds {
		// 检查ip地址判断是否回环地址
		if ip, ok := address.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				host = ip.IP.String()
			}

		}
	}
	return host
}

func main() {
	//注册参数
	flag.StringVar(&phpPath, "php", "php", "PHP位置")
	flag.IntVar(&logger, "log", 0, "是否打印日志,1打印")
	flag.StringVar(&ipv4, "ip", getClientIp(), "服务IP地址，默认本地IP")
	flag.StringVar(&port, "port", "0", "服务端口，默认随机")
	flag.Parse()
	//启动服务
	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", ipv4, port))
	if err != nil {
		panic(err)
	}
	lis, err := net.ListenTCP("tcp", address)
	if err != nil {
		panic(err)
	}
	defer func(lis *net.TCPListener) {
		err := lis.Close()
		if err != nil {
			panic(err)
		}
	}(lis)

	//写入配置文件
	if err := ioutil.WriteFile("proxy.php", []byte("<?php return [\"host\"=>\""+lis.Addr().String()+"\"];"), 0644); err != nil {
		panic(err)
	}

	// 实例化grpc服务端
	s := grpc.NewServer()

	// 注册Greeter服务
	RemoteProxy.RegisterConsumerServer(s, &server{})

	// 往grpc服务端注册反射服务
	reflection.Register(s)

	log.Printf("\n+---------KYY-PROXY---------\n| [host  ] %s\n| [server] running...\n+---------------------------", lis.Addr().String())

	// 启动grpc服务
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
