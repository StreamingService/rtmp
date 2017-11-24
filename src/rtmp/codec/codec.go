/*编码、解码*/

package codec

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