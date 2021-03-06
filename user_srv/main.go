package main

import (
	"flag"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"shop_srvs/user_srv/handler"
	"shop_srvs/user_srv/proto"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50051, "端口号")

	flag.Parse()
	fmt.Println("ip: ", *IP)
	fmt.Println("port: ", *Port)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	err = server.Serve(listen)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}
