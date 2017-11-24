package msg

type SetChunkSize struct {
	Header Header
	ChunkSize uint32 // 4byte
}