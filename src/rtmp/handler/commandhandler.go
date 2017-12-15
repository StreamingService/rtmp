package handler

import (
	"bytes"
	"errors"
	"log"
	"strconv"
	
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

	switch cmd.Name {
	case "connect":
		return handleConnect(c, cmd)
	case "_checkbw":
		return handle_checkbw(c, cmd)
	case "releaseStream":
		return handleReleaseStream(c, cmd)
	case "FCPublish":
		return handleFCPublish(c, cmd)
	case "createStream":
		return handleCreateStream(c, cmd)
	case "publish":
		return handlePublish(c, cmd)
	}

	return nil
}

/*
处理连接命令
*/
func handleConnect(c client.Client, cmd *msg.Command) error {
	log.Print("开始处理connect命令")
	app , isStr := cmd.Object["app"].(string)
	if (!isStr) {
		return errors.New("无app字段")
	}

	log.Printf("客户端连接应用: %s", app)

	// ！！！不知道这段意义,但要加了这些数据, 客户端在发送publish命令后才发继续发送推流数据
	// 本段数据从fms4.5网络包中抓取
	unknownBytes := []byte{ 
		0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x26, 0x25, 0xa0,
		0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x26, 0x25, 0xa0,
		0x02, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04,
		0x00, 
	}

	buf := bytes.NewBuffer(unknownBytes)
	
	// 给客户端发送_result命令消息
	resultCmd := msg.Command {
		Header: msg.Header {
			Format: 0,
			ChunkStreamId: cmd.Header.ChunkStreamId,
		},
		Name: "_result",
		TransactionId: cmd.TransactionId,
		Object: map[string]interface{} {
			"fmsVer": "FMS/4,5,0,297",
			"capabilities": float64(255),
			"mode": float64(1),
		},
		UserArguments: []interface{} {
			map[string]interface{} {
				"level": "status",
				"code": "NetConnection.Connect.Success",
				"description": "Connection succeeded.",
				"objectEncoding": float64(0),
				"data": []interface{} { // 数组
					map[string]interface{} { // 数组中只支持obj序列化
						"version": "4,5,0,297",
					},
				},
			},
		},
	}
	buf.Write(resultCmd.Encode())

	_, err := c.Write(buf.Bytes())
	return err	
}

/*
处理 _checkbw
*/
func handle_checkbw(c client.Client, cmd *msg.Command) error {
	log.Print("开始处理_checkbw命令")


	// 触发onFCPublish事件
	onFCPublish := msg.Command {
		Header: msg.Header {
			Format: 1,
			ChunkStreamId: cmd.Header.ChunkStreamId,
		},
		Name: "onFCPublish",
		TransactionId: 0,
		Object: nil,
		UserArguments: []interface{} {
			map[string]interface{} {
				"code": "NetStream.Publish.Start",
				"description": "description2",
				"details": "details2",
				"clientid": strconv.FormatInt(int64(c.GetIndex()), 10),
			},
		},
	}
	c.Write(onFCPublish.Encode())

	// result
	_result := msg.Command {
		Header: msg.Header {
			Format: 1,
			ChunkStreamId: cmd.Header.ChunkStreamId,
		},
		Name: "_result",
		TransactionId: cmd.TransactionId,
	}

	_, err := c.Write(_result.Encode())

	return err
}

func handleReleaseStream(c client.Client, cmd *msg.Command) error {
	log.Print("开始处理releaseStream命令")


	// onBWDone
	onBWDone := msg.Command {
		Header: msg.Header {
			Format: 1,
			ChunkStreamId: cmd.Header.ChunkStreamId,
		},
		Name: "onBWDone",
		TransactionId: 0,
	}
	c.Write(onBWDone.Encode())

	// 给客户端发送_result命令消息
	respCmd := msg.Command {
		Header: msg.Header {
			Format: 0,
			ChunkStreamId: cmd.Header.ChunkStreamId,
		},
		Name: "_result",
		TransactionId: cmd.TransactionId,
	}

	_, err := c.Write(respCmd.Encode())

	return err
}


func handleFCPublish(c client.Client, cmd *msg.Command) error {
	// 触发onFCPublish事件
	onFCPublish := msg.Command {
		Header: msg.Header {
			Format: 1,
			ChunkStreamId: cmd.Header.ChunkStreamId,
		},
		Name: "onFCPublish",
		TransactionId: cmd.TransactionId,
		Object: nil,
		UserArguments: []interface{} {
			map[string]interface{} {
				"code": "NetStream.Publish.Start",
				"description": "description2",
				"details": "details2",
				"clientid": strconv.FormatInt(int64(c.GetIndex()), 10),
			},
		},
	}
	_, err := c.Write(onFCPublish.Encode())
	return err
}

func handleCreateStream(c client.Client, cmd *msg.Command) error {
	log.Print("处理createStream命令")
	_result := msg.Command {
		Header: msg.Header {
			Format: 0,
			ChunkStreamId: cmd.Header.ChunkStreamId,
		},
		Name: "_result",
		TransactionId: cmd.TransactionId,
		UserArguments: []interface{} {
			float64(1), // stream id
		},
	}
	_, err := c.Write(_result.Encode())

	return err
}

func handlePublish(c client.Client, cmd *msg.Command) error {
	log.Print("处理publish命令")

	// c.UserArguments[0] : publish name
	// c.UserArguments[1] : publish type, live、record、append
	if (cmd.UserArguments == nil || len(cmd.UserArguments) < 2) {
		log.Print("publish user argument length < 2")
		return nil
	}

	publishName, isStr := cmd.UserArguments[0].(string)
	if (!isStr) {
		log.Printf("name数据不是string类型")
		return nil
	}

	// 触发onStatus事件
	onStatus := msg.Command {
		Header: msg.Header {
			Format: 0,
			ChunkStreamId: cmd.Header.ChunkStreamId,
			StreamId: cmd.Header.StreamId,
		},
		Name: "onStatus",
		TransactionId: cmd.TransactionId,
		Object: nil,
		UserArguments: []interface{} {
			map[string]interface{} {
				"level": "status",
				"code": "NetStream.Publish.Start",
				"description": publishName + " is now published.",
				"details": publishName,
				"clientid": strconv.FormatInt(int64(c.GetIndex()), 10),
			},
		},
	}
	_, err := c.Write(onStatus.Encode())
	return err
}