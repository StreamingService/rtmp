package handler

import (
	"log"
	
	"client"
	"rtmp/msg"
	"rtmp/session"
)

type SetChunkSizeHandler struct {

}

func (h *SetChunkSizeHandler) Handle(se *session.Session, c client.Client, m msg.ClientMsg) error {
	scsMsg, ok := m.(*msg.SetChunkSize) // 消息类型转换
	if (ok) {
		log.Print("处理ChunkSizeHandler消息")
		c.MaxChunkSize = scsMsg.ChunkSize
	}

	return nil
}