package handler

import (
	"log"
	"reflect"

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

var handlerMap map[string]MsgHandler = make(map[string]MsgHandler)

/*
获取消息的处理器
*/
func GetMsgHandler(m msg.ClientMsg) MsgHandler {
	msgName := reflect.TypeOf(m).String() // 获取消息实例名称
	return handlerMap[msgName]
}

/*
注册消息处理器
*/
func RegistHandlers() {
	registHandler("*msg.SetChunkSize", &SetChunkSizeHandler{})
	registHandler("*msg.Command", &CommandHandler{})
}

func registHandler(msgType string, handler MsgHandler) {
	log.Printf("注册消息处理器: %s = %s", msgType, reflect.TypeOf(handler).String())
	handlerMap[msgType] = handler
}