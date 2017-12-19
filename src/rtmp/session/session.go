package session


/*
服务会话数据
*/
type Session struct {
	data map[string]interface{}
}

func NewSession() *Session {
	return &Session {
		data: map[string]interface{} {
			"MaxChunkSize": uint32(128),
		},
	}
}

/*
获取属性
*/
func (s *Session) GetAttr(name string) interface{} {
	return s.data[name]
}

/*
设置属性
*/
func (s *Session) SetAttr(name string, value interface{}) {
	s.data[name] = value
}