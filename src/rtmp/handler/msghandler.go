package handler

import (
	"client"
	"rtmp/msg"
)

/*
消息处理器 接口
*/
type MsgHandler interface {

	/* 处理消息 */
	Handle(c client.Client, m msg.ClientMsg) error
}

/*
获取消息的处理器
*/
func GetMsgHandler(m msg.ClientMsg) MsgHandler {
	return nil
}