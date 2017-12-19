package handler

import (
	"rtmp/session"
	"rtmp/msg"
	"client"
)

type AudioHandler struct {

}

func (h *AudioHandler) Handle(se *session.Session, c client.Client, m msg.ClientMsg) error {
	audio, isAudio := m.(*msg.Audio)
	if (!isAudio) {
		return nil
	}

	// 写audio数据到流
	return write2streaming(se, audio.Data)
}