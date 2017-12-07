package service

import (
	"log"
	"errors"

	"rtmp/msg"
)

/*
握手
*/
func (s *RtmpService) handshake() error {
	log.Print("开始handshake")

	c0 := msg.C0 {}
	err := c0.Read(&s.client)
	if (err != nil) {
		return err
	}
	log.Printf("客户端协议版本:%d", c0.Version)
	if (c0.Version != 3) {
		return errors.New("客户端协议不为3,不支持的协议版本")
	}

	// C1
	c1 := msg.C1 {}
	err = c1.Read(&s.client)
	if (err != nil) {
		return err
	}
	log.Printf("c1时间:%d", c1.GetTime())

	// S0
	s0 := msg.S0 {
		Version: 3, // 版本号3
	}
	err = s.sendMsg(&s0)
	if (err != nil) {
		return err
	}

	// S1
	s1 := msg.S1 {
		Time: 10,
	}
	err = s.sendMsg(&s1)
	if (err != nil) {
		return err
	}

	// C2
	c2 := msg.C2 {}
	err = c2.Read(&s.client)
	if (err != nil) {
		return err
	}
	log.Printf("c2: %#v", c2)

	// S2
	s2 := msg.S2 {
		TimeInC1: c1.GetTime(),
		TimeSendbyS1: s1.Time,
		RandomInC1: c1.Random,
	}
	s.sendMsg(&s2)

	log.Printf("完成handshake")

	return nil
}