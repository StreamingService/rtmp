package handler

import (
	"rtmp/session"
	"rtmp/msg"
	"client"
)

type AudioHandler struct {

}

func (h *AudioHandler) Handle(se session.Session, c client.Client, m msg.ClientMsg) error {
	video, isVideo := m.(*msg.Video)
	if (!isVideo) {
		return nil
	}

	// 写video数据到流
	return write2streaming(se, video.Data)
}