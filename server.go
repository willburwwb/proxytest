package main

import (
	"log"
	"net"
	"proxy/server/proxy"
)

func main() {
	conn, err := net.Listen("tcp", ":1080")
	if err != nil {
		log.Println("sock5 server listen failed")
	}
	log.Println("sock5 server listen 1080")
	defer conn.Close()
	for {
		client, err := conn.Accept()
		if err != nil {
			log.Println("accept client failed")
		}
		go proxy.Handel(client)
	}
}
