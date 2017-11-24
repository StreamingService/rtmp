package service

import (
	"log"

	"client"
	"rtmp/msg"
)

const (
	InitStatus int32 = 0 // 初始状态
	HandshakeDone int32 = 10 // 已建议连接状态
)


/*
rtmp服务
*/
type RtmpService struct {
	client client.Client // 客户端
	status int32 // 状态
}

func New(c client.Client) RtmpService {
	return RtmpService {
		client: c,
		status: InitStatus,
	}
}

/**
执行服务
*/
func (s *RtmpService) DoService() {
	log.Print("执行rtmp服务")

	err := s.handshake()
	if (err != nil) {
		log.Print(err)
		s.client.Close()
	}

	s.status = HandshakeDone // 状态置为已握手
	
	
}

/**
发送服务端消息给客户端
*/
func (s *RtmpService) sendMsg(msg msg.ServerMsg) error {
	writeSize, err := s.client.Write(msg.Encode())
	log.Printf("发送数据到客户端, size:%d", writeSize)
	return err
}