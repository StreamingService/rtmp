package msg

import (
	"io"
	"bytes"
	"log"

	"rtmp/session"
)

/*
从流 解析客户端消息, 不包括handshake消息
*/
func ParseClientMsg(reader io.Reader, session *session.Session) (ClientMsg, error) {
	// 先解析header
	header := Header {}
	err := header.Read(reader)
	if (err != nil) {
		return nil, err
	}
	err = processHeader(&header, session)
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

	case 0x12, 0x0F: // AMF0 or AMF3 Data
		msg = &Data {
			Header: header,
		}

	case 0x08:
		msg = &Audio {
			Header: header,
		}

	case 0x09:
		msg = &Video {
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

/*
处理header
*/
func processHeader(header *Header, se *session.Session) error {
	if (header.Format == 0) {
		setFormat0Header(se, header.ChunkStreamId, header)
		// 重新设置format为0的header到session
	} else {
		// 从session中获取format为0的header, 复制数据到当前header
		format0Header := getFormat0Header(se, header.ChunkStreamId)
		if (format0Header != nil) {
			// 复制数据
			if (header.Format == 1) {
				header.StreamId = format0Header.StreamId
			} else if (header.Format == 2) {
				header.StreamId = format0Header.StreamId
				header.BodySize = format0Header.BodySize
				header.TypeId = format0Header.TypeId
			} else if (header.Format == 3) {
				header.StreamId = format0Header.StreamId
				header.BodySize = format0Header.BodySize
				header.TypeId = format0Header.TypeId
				header.Timestamp = format0Header.Timestamp
			}
		}
	}

	return nil
}

const format0headerMapSessionKey = "msg.format0headerMap"

/*
从session中获取format为0的header
*/
func getFormat0Header(se *session.Session, chunkStreamId uint32) *Header {
	headerMap := se.GetAttr(format0headerMapSessionKey)
	if (headerMap == nil) {
		return nil
	}

	headerMap2, isTypeOk := headerMap.(map[uint32]*Header)
	if (!isTypeOk) {
		log.Printf("!session[%s]数据类型不正确", format0headerMapSessionKey)
		return nil
	}

	return headerMap2[chunkStreamId]
}

/*
设置format为0的header到session
*/
func setFormat0Header(se *session.Session, chunkStreamId uint32, header *Header) {
	var headerMap2 map[uint32]*Header
	headerMap := se.GetAttr(format0headerMapSessionKey)
	if (headerMap == nil) {
		headerMap2 = map[uint32]*Header {}
		se.SetAttr(format0headerMapSessionKey, headerMap2)
	} else {
		headerMap3, isTypeOk := headerMap.(map[uint32]*Header)
		if (isTypeOk) {
			headerMap2 = headerMap3
		} else {
			headerMap2 = map[uint32]*Header {}
			se.SetAttr(format0headerMapSessionKey, headerMap2)
		}
	}

	headerMap2[chunkStreamId] = header
}