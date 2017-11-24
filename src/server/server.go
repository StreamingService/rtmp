package server

import (
	"net"
	"log"

	"client"
	"rtmp/service"
)

type Server struct {
	status int8
	bind string
	socket net.Listener
}

const NewStatus int8 = 0
const InitStatus int8 = 1


func New(bind string) (Server, error) {
	log.Printf("new serv %s", bind)
	server := Server {
		status: NewStatus,
		bind: bind,
	}

	socket, err := net.Listen("tcp", bind)
	if (err != nil) {
		return server, err
	}

	server.socket = socket

	return server, nil
}

/*
启动服务
*/
func (s *Server) Serv() (error) {
	log.Printf("accept")

	var clientIndex int32 = 0
	
	for {
		clientConn, err := s.socket.Accept()
		if (err != nil) {
			return err
		}

		clientIndex++
		c := client.New(clientConn, clientIndex)
		
		go servClient(c) // 开启协程
	}

}

/*
为客户端提供服务
*/
func servClient(c client.Client) {
	rtmpService := service.New(c)
	rtmpService.DoService()
}

/*
关闭服务
*/
func (s *Server) Close() {
	s.socket.Close()
}