package client

import (
	"net"
	"log"
	"time"

	//"rtmp/parser"
)

/**
客户端信息
*/
type Client struct {
	clientConn net.Conn
	index int32
	createTime time.Time
	version int8 // 协议版本
}

func New(clientConn net.Conn, index int32) Client {
	defer log.Printf("创建客户端:%d", index)
	return Client {
		clientConn: clientConn,
		index: index,
		createTime: time.Now(),
	}
}

func (c *Client) Read(b []byte) (n int, err error) {
	return c.clientConn.Read(b)
}

func (c *Client) Write(b []byte) (n int, err error) {
	return c.clientConn.Write(b)
}

func (c *Client) Close() error {
	return c.clientConn.Close()
}