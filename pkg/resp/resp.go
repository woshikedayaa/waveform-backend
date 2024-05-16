package resp

import "time"

type Resp struct {
	Code      int    `json:"code"`
	Message   any    `json:"message"`
	IsError   bool   `json:"isError"`
	Err       string `json:"err"`
	Timestamp int64  `json:"ts"`
}

func New(code int, msg any, IsErr bool, Err string) Resp {
	return Resp{
		Code:      code,
		Message:   msg,
		IsError:   IsErr,
		Err:       Err,
		Timestamp: time.Now().Unix(),
	}
}

func Success(obj any) Resp {
	return New(200, obj, false, "")
}

func Fail(err string) Resp {
	return New(400, "", true, err)
}

func Error(err string) Resp {
	return New(500, "", true, err)
}
