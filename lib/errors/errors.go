package errors

import "encoding/json"

type ErrCode string

type Error struct {
	Code ErrCode
	Msg  string
	meta map[interface{}]interface{}
}

func (e *Error) Error() string {
	v, _ := json.Marshal(e)
	return string(v)
}

func New(code ErrCode, err error, keyVals ...interface{}) *Error {
	meta := make(map[interface{}]interface{})
	for i := 0; i < len(keyVals); i += 2 {
		meta[keyVals[i]] = keyVals[i+1]
	}

	return &Error{
		Code: code,
		Msg:  err.Error(),
		meta: meta,
	}
}
