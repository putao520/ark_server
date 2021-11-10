package common

import (
	"encoding/json"
	"server/trans"
)

type ResponseResult struct {
	Code trans.ResultCode
	Message string
	Data interface{}
}

func ResponseResultNew(code trans.ResultCode, message string, data interface{}) []byte{
	result, _ := json.Marshal(&ResponseResult{
		Code:    code,
		Message: message,
		Data:    data,
	})
	return result
}