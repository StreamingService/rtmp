package msg

import (
	"bytes"
	"math/rand"

	"rtmp/codec"
)


/**
服务端消息S1
length 1536 bit
*/
type S1 struct {
	Time uint32
}

func (msg *S1) Encode() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 1536))
	buf.Write(codec.EnInt32(msg.Time)) // 32 bit time
	buf.Write(codec.EnInt32(0)) // 32 bit 0

	// 1528 bit rondom data
	dataSize, _ := buf.WriteString("rtmp streaming server; version: 0.0.1; author: zhujun; ")
	for i:=0; i < 1528 - dataSize; i++  {
		buf.WriteByte(byte(rand.Uint32()))
	}

	return buf.Bytes()
}