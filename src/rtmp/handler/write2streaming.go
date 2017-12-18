package handler

import (
	"errors"

	"rtmp/streaming"
	"rtmp/session"
)

/*
写数据到流
*/
func write2streaming(se session.Session, data []byte) error {
	st := se.GetAttr("msg.streaming")
	if (st == nil) {
		return errors.New("session中无streaming实例")
	}

	st2, isTypeOk := st.(streaming.Streaming)
	if (!isTypeOk) {
		return errors.New("session中streaming类型不正确")
	}

	return st2.Write(data)
}