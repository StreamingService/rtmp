package msg

import (
	"io"
	"errors"

	"rtmp/codec"
)

type C2 struct {
	TimeInS1 uint32 // 32bit
	TimeSendByC1 uint32 // 32bit
}

func (msg *C2) decode(b []byte) {
	msg.TimeInS1 = codec.DeInt32(b)
	msg.TimeSendByC1 = codec.DeInt32(b[4:8])
}

func (msg *C2) Read(reader io.Reader) error {
	bytes := make([]byte, 1536)
	readSize, err := reader.Read(bytes)

	if (err != nil) {
		return err
	}

	if (readSize != 1536) {
		return errors.New("读C2消息长度不足1536")
	}

	msg.decode(bytes)
	return nil
}