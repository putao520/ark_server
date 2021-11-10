package trans

type RealMessage struct {
	Code ResultCode
	Cmd  resultCmd
	Content interface{}
}

func RealMessageNew(code ResultCode, cmd resultCmd, content interface{}) *RealMessage{
	return &RealMessage{
		Code: code,
		Cmd: cmd,
		Content: content,
	}
}

