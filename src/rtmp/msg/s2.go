package msg

import (
	"bytes"
	"math/rand"

	"rtmp/codec"
)

type S2 struct {
	TimeInC1 uint32 // 32bit
	TimeSendbyS1 uint32 // 32bit
}

func (msg *S2) Encode() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 1536))
	buf.Write(codec.EnInt32(msg.TimeInC1)) // 32 bit time
	buf.Write(codec.EnInt32(msg.TimeSendbyS1)) // 32 bit time2

	// 1528 bit rondom data
	for i:=0; i < 1528 ; i++  {
		buf.WriteByte(byte(rand.Uint32()))
	}

	return buf.Bytes()
}