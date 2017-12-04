package msg

import (
	"errors"
	"io"
	
	"rtmp/codec"
)

/*
chunk header
*/
type Header struct {
	/*
	0, 包括所以字段信息, 流的第一个消息format必需为0;
	1, 不包括message stream id, 继承之前消息的;
	2, 不包括message stream id、body size, 继承之前消息的; 用于固定长度的message开始
	3, 不包括message header; 类型0与3组合用于body size > max window zise时, 消息拆分
	*/
	Format uint8 // 2bit, 0、1、2、3
	ChunkStreamId uint32 // 6bit or 14bit or 22bit
	Timestamp uint32 // 4byte
	BodySize uint32 // 3byte
	TypeId uint8 // 1byte
	StreamId uint32 // 4byte
}

/*
从流中解析header
*/
func (h *Header) Read(reader io.Reader) error {
	// begin base header
	// format
	bit1 := []byte{ 0 }
	_, err := reader.Read(bit1)
	if (err != nil) {
		return err
	}
	h.Format = bit1[0] >> 6

	chunkStreamIdType := bit1[0] & 0x3F // 取后6bit
	if (chunkStreamIdType == 0) {
		// 14bit ChunkStreamId, 需要再读1byte
		bit2 := []byte{ 0 }
		_, err = reader.Read(bit2)
		h.ChunkStreamId = uint32(bit2[0]) + 64 // 64 - 319
	} else if (chunkStreamIdType == 1) {
		// 22bit ChunkStreamId, 需要再读2byte
		bit2_3 := []byte { 0, 0 }
		_, err = reader.Read(bit2_3)
		h.ChunkStreamId = uint32(bit2_3[1]) << 8 + uint32(bit2_3[0]) + 64 // 64 - 65599
	} else if (chunkStreamIdType > 1 && chunkStreamIdType < 64) { // 2 - 63, 2为保留值
		h.ChunkStreamId = uint32(chunkStreamIdType)
	}

	// begin message header
	if (h.Format == 0) {
		err = h.readType0MsgHeader(reader)
	} else if (h.Format == 1) {
		err = h.readType1MsgHeader(reader)
	} else if (h.Format == 2) {
		err = h.readType2MsgHeader(reader)
	} else if (h.Format == 3) {
		err = h.readType3MsgHeader(reader)
	}
	if (err != nil) {
		return err
	}

	// extended timestamp
	if (h.Timestamp == 0xFFFFFF) {
		h.Timestamp, err = read4byteInt(reader)
		if (err != nil) {
			return err
		}
	}

	return nil
}

/*
Type0 Message Header
*/
func (h *Header) readType0MsgHeader(reader io.Reader) error {
	time, err := read3byteInt(reader)
	if (err != nil) {
		return err
	}
	h.Timestamp = time

	bodySize, err := read3byteInt(reader)
	if (err != nil) {
		return err
	}
	h.BodySize = bodySize

	typeId, err := readByte(reader)
	if (err != nil) {
		return err
	}
	h.TypeId = typeId

	streamId, err := read4byteInt(reader)
	if (err != nil) {
		return err
	}
	h.StreamId = streamId

	return nil
}

func (h *Header) readType1MsgHeader(reader io.Reader) error {
	time, err := read3byteInt(reader)
	if (err != nil) {
		return err
	}
	h.Timestamp = time

	bodySize, err := read3byteInt(reader)
	if (err != nil) {
		return err
	}
	h.BodySize = bodySize

	typeId, err := readByte(reader)
	if (err != nil) {
		return err
	}
	h.TypeId = typeId

	return nil
}

func (h *Header) readType2MsgHeader(reader io.Reader) error {
	time, err := read3byteInt(reader)
	if (err != nil) {
		return err
	}
	h.Timestamp = time

	return nil
}

func (h *Header) readType3MsgHeader(reader io.Reader) error {
	// no msg header

	return nil
}


func read3byteInt(reader io.Reader) (uint32, error) {
	b := make([]byte, 3)
	_, err := reader.Read(b)
	if (err != nil) {
		return 0, err
	}

	return codec.DeInt24(b), nil
}

func read4byteInt(reader io.Reader) (uint32, error) {
	b := make([]byte, 4)
	_, err := reader.Read(b)
	if (err != nil) {
		return 0, err
	}

	return codec.DeInt32(b), nil
}

func readByte(reader io.Reader) (byte, error) {
	b := []byte{ 0 }
	_, err := reader.Read(b)
	return b[0], err
}

/*
将header数据序列化到输出流
*/
func (h *Header) Write(writer io.Writer) error {
	err := h.writeBaseHeader(writer)
	if (err != nil) {
		return err
	}

	err = h.writeMsgHeader(writer)
	if (err != nil) {
		return err
	}

	return nil
}

/*
写 formatid 和 chunk stream id
*/
func (h *Header) writeBaseHeader(writer io.Writer) error {
	var err error
	if (h.ChunkStreamId < 64) {
		// base header为1 byte
		_, err = writer.Write([]byte{ h.Format << 6 | uint8(h.ChunkStreamId) })
	} else if (h.ChunkStreamId < 320) {
		// base header为2 byte
		_, err = writer.Write([]byte{ h.Format << 6, uint8(h.ChunkStreamId - 64) })
	} else if (h.ChunkStreamId < 65600) {
		// base header为3 byte
		temp := h.ChunkStreamId - 64
		_, err = writer.Write([]byte{ h.Format << 6 | 1,  uint8(temp >> 8), uint8(temp)})
	} else {
		err = errors.New("ChunkStreamId超过支持最大值65599")
	}
	return err
}

/* 
写chunk message header
*/
func (h *Header) writeMsgHeader(writer io.Writer) error {
	switch h.Format {
	case 0:
		return h.writeType0MsgHeader(writer)
	case 1:
		return h.writeType1MsgHeader(writer)
	case 2:
		return h.writeType2MsgHeader(writer)
	case 3:
		return h.writeType3MsgHeader(writer)
	}
	return errors.New("不支持的Format id")
}

func (h *Header) writeType0MsgHeader(writer io.Writer) error {
	writer.Write(codec.EnInt24(h.Timestamp))
	writer.Write(codec.EnInt24(h.BodySize))
	writer.Write([]byte{ h.TypeId })
	writer.Write(codec.EnInt24(h.StreamId))
	return nil
}

func (h *Header) writeType1MsgHeader(writer io.Writer) error {
	writer.Write(codec.EnInt24(h.Timestamp))
	writer.Write(codec.EnInt24(h.BodySize))
	writer.Write([]byte{ h.TypeId })
	return nil
}

func (h *Header) writeType2MsgHeader(writer io.Writer) error {
	writer.Write(codec.EnInt24(h.Timestamp))
	return nil
}

func (h *Header) writeType3MsgHeader(writer io.Writer) error {
	return nil
}
