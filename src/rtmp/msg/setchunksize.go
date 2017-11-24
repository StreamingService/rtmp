package msg

import (
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