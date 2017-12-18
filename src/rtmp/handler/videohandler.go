package handler

import (
	"rtmp/session"
	"rtmp/msg"
	"client"
)

type VideoHandler struct {

}

func (h *VideoHandler) Handle(se session.Session, c client.Client, m msg.ClientMsg) error {
	return nil
}