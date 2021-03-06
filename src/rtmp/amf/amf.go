package amf

import (
	"reflect"
	"log"
	"bytes"
	"strconv"
	"errors"
	"io"

	"rtmp/codec"
)

/*
序列化

Number	0×00	double类型	float64
Boolean	0×01	bool类型		bool
String	0×02	string类型	string
Object	0×03	object类型	map[string]interface{}
MovieClip	0×04	Not available in Remoting
Null	0×05	null类型，空	nil
Undefined	0×06	
Reference	0×07	
MixedArray	0×08				[]map[string]interface{}
EndOfObject	0×09	See Object ，表示object结束
Array	0x0a	
Date	0x0b	
LongString	0x0c	
Unsupported	0x0d	
Recordset	0x0e	Remoting, server-to-client only
XML	0x0f	
TypedObject (Class instance)	0×10	
AMF3 data	0×11	
Sent by Flash player 9+

*/
func Serialize(data interface{}) []byte {
	if (data == nil) {
		return SerializeNull()
	}

	str, isStr := data.(string)
	if (isStr) {
		return SerializeString(str)
	}

	num, isNum := data.(float64)
	if (isNum) {
		return SerializeNumber(num)
	}

	boolean, isBoolean := data.(bool)
	if (isBoolean) {
		return SerializeBoolen(boolean)
	}

	obj, isObj := data.(map[string]interface{})
	if (isObj) {
		return SerializeObject(obj)
	}

	arr, isArr := data.([]map[string]interface{})
	if (isArr) {
		return SerializeArray(arr)
	}

	log.Printf("!!不支持的amf数据类型: %s, 编码为null", reflect.TypeOf(data).String())
	return SerializeNull()
}

/*
反序列化

从流中反序列化一个amf数据
*/
func Deserialize(r io.Reader) (interface{}, error) {
	b1 := []byte {0}
	_, err := r.Read(b1)
	if (err != nil) {
		return nil, err
	}

	switch b1[0] {
	case 0x00:
		return deserializeNumber(r)
	case 0x01:
		return deserializeBoolen(r)
	case 0x02:
		return deserializeString(r)
	case 0x03:
		return deserializeObject(r)
	case 0x05: // Null
		return nil, nil
	case 0x08:
		return deserializeArray(r)
	}

	// 其它类型不支持
	return nil, errors.New("不支持的amf类型: " + strconv.FormatUint(uint64(b1[0]), 16))
}

func SerializeNumber(number float64) []byte {
	buf := bytes.NewBuffer([]byte{ 0x00 }) // 类型
	buf.Write(codec.EnFloat64(number))
	return buf.Bytes()
}

func deserializeNumber(r io.Reader) (float64, error) {
	b := make([]byte, 8)
	_, err := r.Read(b)
	if (err != nil) {
		return 0, err
	}
	return codec.DeFloat64(b), nil
}

func deserializeBoolen(r io.Reader) (bool, error) {
	b := make([]byte, 1)
	_, err := r.Read(b)
	if (err != nil) {
		return false, err
	}
	return b[0] != 0, nil
}

func SerializeBoolen(value bool) []byte {
	var boolByte byte
	if ( value ) {
		boolByte = 1
	} else {
		boolByte = 0
	}
	return []byte{ 0x01, boolByte }
}

func SerializeString(value string) []byte {
	buf := bytes.NewBuffer([]byte{ 0x02 })
	strBytes := []byte(value)
	buf.Write(codec.EnInt16(uint16(len(strBytes))))
	buf.Write(strBytes)
	return buf.Bytes()
}

func deserializeString(r io.Reader) (string, error) {
	b16 := make([]byte, 2) // 长度16bit
	_, err := r.Read(b16)
	if (err != nil) {
		return "", err
	}

	len := codec.DeInt16(b16)
	if (len == 0) {
		return "", nil
	}

	bStr := make([]byte, len)
	_, err = r.Read(bStr)
	if (err != nil) {
		return "", err
	}

	return string(bStr), nil
}

func deserializeObject(r io.Reader) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	b3 := make([]byte, 3) // 用于end测试
	for {
		_, err := r.Read(b3)
		if (err != nil) {
			return obj, err
		}

		if (isObjEnd(b3)) {
			// 到达对象结束标识
			break
		}

		// 未结束, b3的前2byte为字段名称长度
		proLen := codec.DeInt16(b3)
		proAfter := make([]byte, proLen - 1) // 长度需要-1, 之前b3已经多读了一个字节
		_, err = r.Read(proAfter) // 读字段名
		if (err != nil) {
			return obj, err
		}

		// 组装完整proName []byte
		proNameBuf := bytes.NewBuffer([]byte{b3[2]})
		proNameBuf.Write(proAfter)
		proName := string(proNameBuf.Bytes())

		value, err := Deserialize(r) // 反序列化值, 这里有递归调用
		obj[proName] = value
	}

	return obj, nil
}

func isObjEnd(b3 []byte) bool {
	if (len(b3) < 3) {
		return false
	}

	return b3[0] == 0x00 && b3[1] == 0x00 && b3[2] == 0x09
}

/*
amf序列化obj
*/
func SerializeObject(obj map[string]interface{}) []byte {
	if (obj == nil) {
		return SerializeNull()
	}

	buf := bytes.NewBuffer([]byte{ 0x03 }) // 对象类型头
	keyValBytes, _ := serializeObjectKeyValue(obj)
	buf.Write(keyValBytes)
	buf.Write([]byte{ 0x00, 0x00, 0x09 }) // 对象结束

	return buf.Bytes()
}

/*
序列化obj, 不包括类型和结尾
返回 编码数据和项目数
*/
func serializeObjectKeyValue(obj map[string]interface{}) ([]byte, uint32) {
	buf := bytes.NewBuffer([]byte{ })
	var itemSize uint32 = 0
	for k, v := range obj {
		itemSize++
		// 序列化属性名
		strBytes := []byte(k)
		buf.Write(codec.EnInt16(uint16(len(strBytes))))
		buf.Write(strBytes)

		// 序列化值
		buf.Write(Serialize(v))
	}
	return buf.Bytes(), itemSize
}


func SerializeNull() []byte {
	return []byte{ 0x05 }
}

/*
序列化数组
*/
func SerializeArray(arr []map[string]interface{}) []byte {
	if (arr == nil) {
		return SerializeNull()
	}

	buf := bytes.NewBuffer([]byte{ 0x08 })
	var arrLen uint32 = 0
	objBuf := bytes.NewBuffer([]byte{})
	for _, v := range arr {
		// 只序列化数组中的对象
		keyValBytes, itemSize := serializeObjectKeyValue(v)
		arrLen += itemSize
		objBuf.Write(keyValBytes);
	}
	buf.Write(codec.EnInt32(arrLen)) // 长度
	buf.Write(objBuf.Bytes()) // 数据
	buf.Write([]byte{ 0x00, 0x00, 0x09 }) // 结尾

	return buf.Bytes()
}

func deserializeArray(r io.Reader) ([]map[string]interface{}, error) {
	lenBytes := make([]byte, 4)
	_, err := r.Read(lenBytes)
	if (err != nil) {
		return nil, err
	}
	log.Printf("读取array数据项length: %d", codec.DeInt32(lenBytes))
	obj, err := deserializeObject(r) // 解析字段

	return []map[string]interface{}{ obj }, nil
}