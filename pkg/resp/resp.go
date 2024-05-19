// 响应封装

package resp

import "time"

// 封装 HTTP 响应信息
type Resp struct {
	Code      int    `json:"code"`    //响应状态码
	Message   any    `json:"message"` //响应数据
	IsError   bool   `json:"isError"` //是否为错误响应
	Err       string `json:"err"`     //错误信息，当“IsError”为true时起作用
	Timestamp int64  `json:"ts"`      //响应时间戳
}

// 根据传入参数设置响应
func New(code int, msg any, IsErr bool, Err string) Resp {
	return Resp{
		Code:      code,
		Message:   msg,
		IsError:   IsErr,
		Err:       Err,
		Timestamp: time.Now().Unix(),
	}
}

// 成功响应封装
func Success(obj any) Resp {
	return New(200, obj, false, "")
}

// 失败响应封装
func Fail(err string) Resp {
	return New(400, "", true, err)
}

// 错误响应封装
func Error(err string) Resp {
	return New(500, "", true, err)
}
