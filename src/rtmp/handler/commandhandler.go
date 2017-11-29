package handler

import (
	"errors"
	"log"

	"client"
	"rtmp/msg"
)

/*
命令处理
*/
type CommandHandler struct {

}

func (h *CommandHandler) Handle(c client.Client, m msg.ClientMsg) error {
	log.Print("处理Command消息")
	cmd, ok := m.(*msg.Command) // 消息类型转换
	if (!ok) {
		return errors.New("消息类型不正确")
	}

	if (cmd.Name == "connect") {
		return handleConnect(c, cmd)
	}

	return nil
}

/*
处理连接命令
*/
func handleConnect(c client.Client, cmd *msg.Command) error {
	app , isStr := cmd.Object["app"].(string)
	if (!isStr) {
		return errors.New("无app字段")
	}

	log.Printf("客户端连接应用: %s", app)

	// 给客户端发送onBWDone命令消息

	onBWDone := msg.Command {
		Header: msg.Header {
			ChunkStreamId: cmd.Header.ChunkStreamId,
		},
		Name: "onBWDone",
		TransactionId: 0,
	}

	_, err := c.Write(onBWDone.Encode())

	return err
}