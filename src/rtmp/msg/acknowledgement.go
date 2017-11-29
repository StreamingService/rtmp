package msg

/*
Acknowledgement

typeid: 0x03

当接收到的消息大小为 window size时, 接收方需要回复发送方该消息

*/
type Acknowledgement struct {
	SequenceNumber uint32 // the number of bytes received so far
}