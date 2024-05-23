package wave

import "strings"

type OpError struct {
	err        error
	op         string
	port       string
	suggestion string
}

func (o *OpError) Error() string {
	// 可能会有疑问这里为什么要判定nil
	// 因为golang可以这样调用对应的方法
	/*
		oe := &OpError{}
		OpError.Error(oe)
		// 和下面这个是等价的
		oe.Error()
		// 但是如果 oe 变成了 nil
		OpError.Error(nil)
		...
	*/
	if o == nil {
		return "<nil>"
	}
	es := ""
	if o.err != nil {
		es = o.err.Error()
	}

	return strings.Join([]string{"waveform: ", o.op, o.port, es, o.suggestion}, " ")
}
