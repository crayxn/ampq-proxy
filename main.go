package main

import (
	AmqpProxy "amqp-proxy/proto"
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
var artisanPath string
var artisanCommand string
var logger int

func logConsole(format string, v ...any) {
	if logger == 1 {
		log.Printf(format, v...)
	}
}

func (s *server) Do(ctx context.Context, params *AmqpProxy.Params) (*AmqpProxy.Reply, error) {

	logConsole("%s %s %s %s %s %s", phpPath, artisanPath, artisanCommand, params.GetClass(), params.GetFunc(), params.GetArgs())
	cmd := exec.Command(phpPath, artisanPath, artisanCommand, params.GetClass(), params.GetFunc(), "--args="+params.GetArgs())
	out, err := cmd.CombinedOutput()
	logConsole("combined out:\n%s\n", string(out))
	if err != nil {
		logConsole("combined out:\n%s\n", string(out))
		return &AmqpProxy.Reply{Success: false, Message: fmt.Sprintf("执行错误,请检查命令 %s -> %s(%s)", params.GetClass(), params.GetFunc(), params.GetArgs())}, nil
	}
	if string(out) != "success" {
		return &AmqpProxy.Reply{Success: false, Message: string(out)}, nil
	}
	return &AmqpProxy.Reply{Success: true, Message: ""}, nil
}

func getClientIp() (string, error) {
	adds, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}
	host := "0.0.0.0"
	for _, address := range adds {
		// 检查ip地址判断是否回环地址
		if ip, ok := address.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				host = ip.IP.String()
			}

		}
	}

	return host, nil

}

func main() {
	//注册参数
	flag.StringVar(&phpPath, "php", "php", "php位置，默认环境中php")
	flag.StringVar(&artisanPath, "artisan", "../../../artisan", "artisan位置，默认../../../artisan")
	flag.StringVar(&artisanCommand, "command", "kyy:proxy", "执行命令，默认 kyy:proxy")
	flag.IntVar(&logger, "log", 0, "是否打印日志 1打印")
	flag.Parse()
	//获取闲置端口
	ip, err := getClientIp()
	if err != nil {
		panic(err)
	}
	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", ip))
	if err != nil {
		panic(err)
	}
	lis, err := net.ListenTCP("tcp", address)
	if err != nil {
		panic(err)
	}
	defer lis.Close()

	//写入配置文件
	if err := ioutil.WriteFile("proxy.php", []byte("<?php return [\"host\"=>\""+lis.Addr().String()+"\"];"), 0644); err != nil {
		panic(err)
	}

	// 实例化grpc服务端
	s := grpc.NewServer()

	// 注册Greeter服务
	AmqpProxy.RegisterConsumerServer(s, &server{})

	// 往grpc服务端注册反射服务
	reflection.Register(s)

	// 启动grpc服务
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Println("server running")
}
