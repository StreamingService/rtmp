package streaming

type Streaming interface {
	Close()

	Write(byteArray []byte) (int, error)
}