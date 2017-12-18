package streaming

type Streaming interface {
	Close()

	Write(byteArray []byte) error
}