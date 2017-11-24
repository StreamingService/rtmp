package msg

/*
服务端消息S0
*/
type S0 struct {
	Version uint8
}

func (msg *S0) Encode() []byte {
	return []byte{ msg.Version }
}