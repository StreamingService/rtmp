package msg

import (
	"rtmp/amf"
	"io"
)

/*
数据消息

typeid: 0x12 AMF0 data
typeid: 0x0F AMF3 data
*/
type Data struct {
	Header Header
	Data []interface{}
}

func (msg *Data) Read(reader io.Reader) error {
	for {
		data, err := amf.Deserialize(reader)
		if (err != nil) {
			// 读完数据, 结束
			break;
		}

		if (msg.Data == nil) {
			msg.Data = []interface{} {}
		}
		msg.Data = append(msg.Data, data)
	}

	return nil
}