package msg

import (
	"errors"
	"io"

	"rtmp/amf"
)

/*
AMF0命令

body数据为amf0编码

typeid: 0x14 AMF0
typeid: 0x11 AMF3
*/
type Command struct {
	Header Header
	Name string // command name
	TransactionId uint32
	Object map[string]interface{} // 命令数据对象
	UserExt interface{} // 用户扩展数据, 可选
}

func (msg *Command) Read(reader io.Reader) error {
	 // name
	data, err := amf.Deserialize(reader)
	if (err != nil) {
		return errors.New("未解析到命令名称")
	}

	cmdName, isStr := data.(string)
	if (!isStr) {
		return errors.New("命令名称应该为字符串类型")
	}
	msg.Name = cmdName

	// TransactionId
	data, err = amf.Deserialize(reader)
	if (err != nil) {
		return errors.New("未解析到TransactionId: " + err.Error())
	}

	id, isNumber := data.(float64)
	if (!isNumber) {
		return errors.New("TransactionId应该为Number类型")
	}
	msg.TransactionId = uint32(id)

	// Object
	data, err = amf.Deserialize(reader)
	if (err != nil) {
		return errors.New("未解析到Object")
	}
	if (data != nil) {
		obj, isObj := data.(map[string]interface{})
		if (!isObj) {
			return errors.New("Object类型不正确")
		}
		msg.Object = obj
	}

	// UserExt 可选
	ext, _ := amf.Deserialize(reader)
	msg.UserExt = ext

	return nil
}