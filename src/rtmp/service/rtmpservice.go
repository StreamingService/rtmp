package service

import (
	"log"
	
	"client"
	"rtmp/msg"
	"rtmp/handler"
	"rtmp/session"
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

	// 握手
	err := s.handshake()
	if (err != nil) {
		log.Print(err)
		s.client.Close()
	}

	s.status = HandshakeDone // 状态置为已握手

	// 初始化会话对象
	sess := session.NewSession()
	
	// 服务循环
	for {
		msg, err := msg.ParseClientMsg(&s.client, sess)
		if (err != nil) {
			log.Print(err)
			break
		}

		if (msg != nil) {
			log.Printf("接收到客户端消息: %#v", msg)
			// 处理消息
			h := handler.GetMsgHandler(msg) // 查询处理器
			if (h != nil) {
				// 处理器存在, 执行处理逻辑
				err = h.Handle(sess, s.client, msg)
				if (err != nil) {
					// 处理消息出错
					log.Printf("处理消息出错: %s", err.Error())
					break; // 跳出服务循环
				}
			}
		}

	}

	// 服务结束
	s.client.Close()
}

/**
发送服务端消息给客户端
*/
func (s *RtmpService) sendMsg(msg msg.ServerMsg) error {
	writeSize, err := s.client.Write(msg.Encode())
	log.Printf("发送数据到客户端, size:%d", writeSize)
	return err
}