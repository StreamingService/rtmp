package msg

import (
	"io"
)

/**
rtmp消息
*/
type RTMPMsg interface {

}

/*
服务端发送到客户端消息
*/
type ServerMsg interface {

	/*
	编码
	*/
	Encode() []byte
}

/*
客户端消息
*/
type ClientMsg interface {

	/*
	从流中解析消息
	*/
	Read(io.Reader) error
}