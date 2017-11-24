package msg

import (
	"io"
	"errors"
)

/*
客户端消息C0
*/
type C0 struct {
	Version uint8
}

func (msg *C0) Read(reader io.Reader) error {
	bytes := make([]byte, 1)
	readSize, err := reader.Read(bytes)

	if (err != nil) {
		return err
	}

	if (readSize != 1) {
		return errors.New("读C0消息size为0")
	}

	msg.Version = bytes[0]

	return nil
}
