package msg

import (
	"bytes"
	"io"

	"rtmp/codec"
)

type SetChunkSize struct {
	Header Header
	ChunkSize uint32 // 4byte
}

func (msg *SetChunkSize) Read(reader io.Reader) error {
	// header已由paser统一解析,这里只解析body部分
	b := make([]byte, 4)
	_, err := reader.Read(b)
	if (err != nil) {
		return err
	}

	msg.ChunkSize = codec.DeInt32(b)

	return nil
}

func (msg *SetChunkSize) Encode() []byte {
	msg.Header.BodySize = 4
	msg.Header.TypeId = 0x01

	buf := bytes.NewBuffer([]byte{})
	msg.Header.Write(buf)
	buf.Write(codec.EnInt32(msg.ChunkSize))
	return buf.Bytes()
}