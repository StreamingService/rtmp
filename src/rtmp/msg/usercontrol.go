package msg

import (
	"rtmp/codec"
	"bytes"
)



/*
用户控制消息

TypeId: 0x04
message stream id: 0
ChunkStreamId: 2
*/
type UserControl struct {
	Header Header
	EventType uint16 // 事件类型
	EventData []byte
}

/*
流开始事件
S->C
默认在收到connect命令后发送给客户端
event data: 4 byte stream id
*/
const UC_StreamBegin uint16 = 0

/*
流结束事件
S->C
event data: 4 byte stream id
*/
const UC_StreamEOF uint16 = 1

/*
通知客户端该流没有更多的数据
S->C
event data: 4 byte stream id
*/
const UC_StreamDry uint16 = 2

/*
C->S
event data: 4 byte stream id, 4 byte buffer length
*/
const UC_SetBufferLength uint16 = 3

/*
通知客户端该stream id已经被录制
S->C
event data: 4 byte stream id
*/
const UC_StreamIsRecorded uint16 = 4

/*
测试客户端是否可达
S->C
event data: 4 byte timestamp
*/
const UC_PingRequest uint16 = 6

/*
客户端响应服务端Ping事件
C->S
event data: 4 byte timestamp 服务端Ping事件中的时间戳
*/
const UC_PingResponse uint16 = 7


func (msg *UserControl) Encode() []byte {
	// 重置header数据
	msg.Header.TypeId = 0x04
	msg.Header.StreamId = 0
	msg.Header.ChunkStreamId = 2

	// body
	bodyBuf := bytes.NewBuffer([]byte{})
	bodyBuf.Write(codec.EnInt16(msg.EventType))
	bodyBuf.Write(msg.EventData)
	body := bodyBuf.Bytes()

	// 重置body size
	msg.Header.BodySize = uint32(len(body))

	msgBuf := bytes.NewBuffer([]byte{})
	msg.Header.Write(msgBuf) // 编码header
	msgBuf.Write(body) // 写body

	return msgBuf.Bytes()
}