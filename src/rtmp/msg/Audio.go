package msg

import (
	"errors"
	"io"
)

/*
type id: 0x08
*/
type Audio struct {
	Header Header
	Format uint8 // 4bit
	SampleRate uint8 // 2bit, 3为44kHz
	SampleSize uint8 // 1bit, 1为16bit
	Channels uint8 // 1bit, 1为stereo(立体声)
	Data []byte // 音频数据, 长度为 body size - 1
}

// begin 音频类型常量
const AF_HEACC uint8 = 10 // he-acc


func (msg *Audio) Read(reader io.Reader) error {
	control := []byte{ 0 }
	_, err := reader.Read(control)
	if (err != nil) {
		return err
	}
	msg.Format = control[0] >> 4
	msg.SampleRate = uint8(control[0] >> 2 & 0x03)
	msg.SampleSize = uint8(control[0] >> 1 & 0x01)
	msg.Channels = uint8(control[0] & 0x01)

	// data
	data := make([]byte, msg.Header.BodySize - 1)
	dataReadSize, err := reader.Read(data)
	if (err != nil) {
		return err
	}
	if (uint32(dataReadSize) != msg.Header.BodySize - 1) {
		return errors.New("Audio data长度不正确")
	}
	msg.Data = data

	return nil
}
