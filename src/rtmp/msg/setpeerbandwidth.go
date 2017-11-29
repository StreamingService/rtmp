package msg

/*
Set Peer Bandwidth

typeid: 0x06

限制对方输出带宽
*/
type SetPeerBandwidth struct {
	Bandwidth uint32 // 带宽
	LimitType uint8 // 限制类型, 0 hard, 1 soft, 2 dynamic
}