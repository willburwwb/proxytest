package proxy

import (
	"log"
	"net"
)

func Handel(conn net.Conn) {
	proxy := &Proxy{
		conn: conn,
	}
	defer conn.Close()
	//认证
	if err := proxy.Auth(); err != nil {
		log.Println("sock5 验证失败")
		return
	}
	log.Println("sock5 验证成功")
	//建立连接
	server, err := proxy.Connect()
	if err != nil {
		log.Println("sock5 连接失败")
		return
	}
	log.Println("sock5 验证失败")
	proxy.server = server
	defer server.Close()

	if err != nil {
		log.Println("connect to server failed")
	}
	//转发数据
	log.Println("开始转发数据")
	proxy.Forward()
}
