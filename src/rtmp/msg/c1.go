package msg

import (
	"errors"
	"io"

	"rtmp/codec"
)

/**
客户端消息C1
length 1536 bit
*/
type C1 struct {
	time uint32
	Random []byte // 随机数据, 1528 byte
}

func (msg *C1) decode(bytes []byte) error {
	if (len(bytes) != 1536) {
		return errors.New("bytes长度应为1536")
	}
	
	// 目前只解析时间字段
	msg.time = codec.DeInt32(bytes)
	msg.Random = bytes[8:]
	return nil
}

func (msg *C1) Read(reader io.Reader) error {
	bytes := make([]byte, 1536)
	readSize, err := reader.Read(bytes)

	if (err != nil) {
		return err
	}

	if (readSize != 1536) {
		return errors.New("读C1消息长度不足1536")
	}

	msg.decode(bytes)
	return nil
}


func (msg *C1) GetTime() uint32 {
	return msg.time
}