package errors

import "encoding/json"

// ErrCode is a type to represent different error codes
type ErrCode string

// Error is a custom error struct that explictly holds
// and ErrCode to differentiate different type of errors
type Error struct {
	Code ErrCode
	Msg  string
	meta map[interface{}]interface{}
}

func (e *Error) Error() string {
	v, _ := json.Marshal(e)
	return string(v)
}

// New returns an error with the given params
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
