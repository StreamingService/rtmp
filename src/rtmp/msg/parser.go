package msg

import (
	"io"
	"bytes"
	"log"
)

/*
从流 解析客户端消息, 不包括handshake消息
*/
func ParseClientMsg(reader io.Reader) (ClientMsg, error) {
	// 先解析header
	header := Header {}
	err := header.Read(reader)
	if (err != nil) {
		return nil, err
	}

	var bodyReader io.Reader
	if (header.BodySize > 0) {
		// 读body
		body := make([]byte, header.BodySize)
		_, err = reader.Read(body)
		if (err != nil) {
			return nil, err
		}
		bodyReader = bytes.NewReader(body)
	}

	var msg ClientMsg

	switch header.TypeId {
	case 0x01 :
		msg = &SetChunkSize {
			Header: header,
		}

	case 0x14 , 0x11: // AMF0 or AMF3 Command
		msg = &Command {
			Header: header,
		}
	}

	if (msg == nil) {
		log.Printf("!不被支持的消息头: %+v", header)
		return nil, nil // 不支持的消息, 返回nil
	}

	// 由消息自已负责解析body
	err = msg.Read(bodyReader)

	return msg, err
} 