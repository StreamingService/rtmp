package msg

import (
	"io"
)

/*
chunk header
*/
type Header struct {
	Format uint8 // 2bit, 0、1、3
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

	return uint32(b[0]) << 16 | uint32(b[1]) << 8 | uint32(b[2]), nil
}

func read4byteInt(reader io.Reader) (uint32, error) {
	b := make([]byte, 4)
	_, err := reader.Read(b)
	if (err != nil) {
		return 0, err
	}

	return uint32(b[0]) << 24 | uint32(b[1]) << 16 | uint32(b[2]) << 8 | uint32(b[3]), nil
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
	return nil
}