package main

import (
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
	token "github.com/light-speak/lightdoor/security/kitex_gen/token/securityservice"
	"github.com/light-speak/lighthouse/env"
)

func main() {
	// 获取服务地址和端口，可以通过环境变量覆盖
	host := env.Getenv("SERVICE_HOST", "0.0.0.0")
	port := env.Getenv("SERVICE_PORT", "4000")

	addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host, port))
	if err != nil {
		log.Fatalf("net.ResolveTCPAddr err: %v", err)
	}
	svr := token.NewServer(new(SecurityServiceImpl), server.WithServiceAddr(addr))
	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
