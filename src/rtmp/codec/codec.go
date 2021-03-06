/*编码、解码*/

package codec

import (
	"math"
	"encoding/binary"
)

func DeInt32(b []byte) uint32 {
	return uint32(b[0]) << 24 | uint32(b[1]) << 16 | uint32(b[2]) << 8 | uint32(b[3])
}

func EnInt32(value uint32) []byte {
	b := make([]byte, 4)
	b[0] = byte(value >> 24)
	b[1] = byte(value >> 16 & 0xFF)
	b[2] = byte(value >> 8 & 0xFF)
	b[3] = byte(value & 0xFF)
	return b
}

func EnInt16(value uint16) []byte {
	b := make([]byte, 2)
	b[0] = byte(value >> 8 & 0xFF)
	b[1] = byte(value & 0xFF)
	return b
}

func DeInt16(b []byte) uint16 {
	return uint16(b[0]) << 8 | uint16(b[1])
}

/*
解码64位浮点数

1位符号位，11位阶码，52位尾数 ??
*/
func DeFloat64(b []byte) float64 {
	i := binary.BigEndian.Uint64(b)
	return math.Float64frombits(i)
}

func EnFloat64(value float64) []byte {
	bits := math.Float64bits(value)
    bytes := make([]byte, 8)
    binary.BigEndian.PutUint64(bytes, bits)
    return bytes
}

func EnInt24(value uint32) []byte {
	b := make([]byte, 3)
	b[0] = byte(value >> 16)
	b[1] = byte(value >> 8 & 0xFF)
	b[2] = byte(value & 0xFF)
	return b
}

func DeInt24(b []byte) uint32 {
	return uint32(b[0]) << 16 | uint32(b[1]) << 8 | uint32(b[2])
}