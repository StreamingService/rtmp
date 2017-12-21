package msg

import (
	"errors"
	"io"
)

/*
type id: 0x09
*/
type Video struct {
	Header Header
	Type uint8 // 4 bit
	Format uint8 // 4 bit
	Data []byte // 视频数据, 长度为body size -1
}

// begin 视频格式常量
const VF_H264 uint8 = 7 // h264 format
const VF_H263 uint8 = 2 // seronson h.263
const VF_ScreenVideo uint8 = 3 // screen video
const VF_On2VP6 uint8 = 4 // On2 VP6
const VF_On2VP6A uint8 = 5 // On2 VP6 with alpha channel
const VF_ScreenVideo2 uint8 = 6 // screen video version 2

// begin 视频类型常量
const VT_Keyframe uint8 = 1 // keyframe
const VT_Innerframe uint8 = 2 // inner frame
const VT_DisposableInnerframe uint8 = 3 // disposable inner frame （h.263 only）
const VT_GeneratedKeyframe uint8 = 4 // generated keyframe

func (msg *Video) Read(reader io.Reader) error {
	control := []byte{ 0 }
	_, err := reader.Read(control)
	if (err != nil) {
		return err
	}
	msg.Type = control[0] >> 4
	msg.Format = control[0] & 0x0F

	// data
	data := make([]byte, msg.Header.BodySize - 1)
	dataReadSize, err := reader.Read(data)
	if (err != nil) {
		return err
	}
	if (uint32(dataReadSize) != msg.Header.BodySize - 1) {
		return errors.New("Video data长度不正确")
	}
	msg.Data = data

	return nil
}