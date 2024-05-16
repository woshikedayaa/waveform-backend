package resp

type Resp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	IsError bool   `json:"isError"`
	Err     string `json:"err"`
}

func New(code int, msg string, IsErr bool, Err string) Resp {
	return Resp{
		Code:    code,
		Message: msg,
		IsError: IsErr,
		Err:     Err,
	}
}

func Success(msg string) Resp {
	return New(200, msg, false, "")
}

func Fail(err string) Resp {
	return New(400, "", true, err)
}

func Error(err string) Resp {
	return New(500, "", true, err)
}
