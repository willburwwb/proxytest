package proxy

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
)

type Proxy struct {
	conn   net.Conn
	server net.Conn
}

func (proxy *Proxy) Auth() (err error) {
	buf := make([]byte, 1024)
	n, err := proxy.conn.Read(buf)
	if err != nil {
		log.Println("read client auth failed")
		return err
	}
	if n <= 1 {
		return errors.New("读取数据失败")
	}
	for _, v := range buf[2:n] {
		log.Printf("client suports method %d\n", v)
	}
	if _, err = proxy.conn.Write([]byte{buf[0], 0x00}); err != nil {
		return err
	}
	return nil
}
func (proxy *Proxy) Connect() (net.Conn, error) {
	buf := make([]byte, 1024)
	n, err := proxy.conn.Read(buf[:4])
	if n != 4 || err != nil {
		log.Println(err)
		if n != 4 {
			return nil, errors.New("读取数据失败")
		}
		return nil, err
	}
	version := buf[0]
	var ip net.IP
	switch buf[3] {
	case 0x01:
		n, err := proxy.conn.Read(buf[:4])
		if n != 4 || err != nil {
			return nil, errors.New("IPv4 invalid" + err.Error())
		}
		//addr = fmt.Sprintf("%s.%s.%s.%s", buf[0], buf[1], buf[2], buf[3])
		ip = buf[0:n]
	case 0x03:
		_, err := proxy.conn.Read(buf[:1])
		if err != nil {
			log.Println(err)
			return nil, err
		}
		length := int(buf[0])
		if _, err = proxy.conn.Read(buf[:length]); err != nil {
			log.Println("读取host err", err)
			return nil, err
		}
		ips, err := net.LookupIP(string(buf[0:length]))
		if err != nil {
			log.Println("from host look up ip failed", string(buf[0:length]))
			return nil, err
		}
		ip = ips[0]
	case 0x04:
		n, err := proxy.conn.Read(buf[:16])
		if n != 16 || err != nil {
			return nil, errors.New("IPv6 invalid" + err.Error())
		}
		ip = buf[0:n]
	}
	_, err = proxy.conn.Read(buf[:2])
	if err != nil {
		return nil, errors.New("Port invalid" + err.Error())
	}

	port := binary.BigEndian.Uint16(buf[0:2])
	addr := net.TCPAddr{
		IP:   ip,
		Port: int(port),
	}
	log.Println("sock5 server connect remote server ", addr)

	server, err := net.Dial("tcp", addr.String())
	if err != nil {
		log.Println("连接 server failed")
		return nil, err
	}
	log.Println("connect remote server successed")
	if _, err := proxy.conn.Write([]byte{version, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0, 0, 0}); err != nil {
		log.Println("send to sock5 client connection message failed")
		return nil, err
	}
	return server, nil
}
func (proxy *Proxy) Transfer(src net.Conn, dst net.Conn) {
	buf := make([]byte, 256)
	for {
		n, err := src.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Println(err)
			return
		}
		dst.Write(buf[0:n])
	}
}
func (proxy *Proxy) Forward() {
	go proxy.Transfer(proxy.server, proxy.conn)
	proxy.Transfer(proxy.conn, proxy.server)
}
