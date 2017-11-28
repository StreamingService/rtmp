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
	log.Print("注册消息处理器: *msg.SetChunkSize=&SetChunkSizeHandler{}")
	handlerMap["*msg.SetChunkSize"] = &SetChunkSizeHandler{}
}