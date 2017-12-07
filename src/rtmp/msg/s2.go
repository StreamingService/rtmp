package msg

import (
	"bytes"

	"rtmp/codec"
)

type S2 struct {
	TimeInC1 uint32 // 32bit
	TimeSendbyS1 uint32 // 32bit
	RandomInC1 []byte // c1中的随机数据, 1528 byte
}

func (msg *S2) Encode() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 1536))
	buf.Write(codec.EnInt32(msg.TimeInC1)) // 32 bit time
	buf.Write(codec.EnInt32(msg.TimeSendbyS1)) // 32 bit time2

	// 1528 bit rondom data from c1
	buf.Write(msg.RandomInC1)

	return buf.Bytes()
}