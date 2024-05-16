package wave

import "strings"

type OpError struct {
	cache string
	raw   error
	op    string
	port  string
	other string
}

func (oe *OpError) Error() string {
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
	if oe == nil {
		return "<nil>"
	}
	if len(oe.cache) != 0 {
		return oe.cache
	}
	if len(oe.port) == 0 {
		panic("wave: OpError must order a port")
	}

	builder := strings.Builder{}
	builder.WriteString("wave: ")
	builder.WriteString(oe.port)
	builder.WriteString(": ")
	if len(oe.op) != 0 {
		builder.WriteString(oe.op)
	}
	if oe.raw != nil {
		builder.WriteString(": ")
		builder.WriteString(oe.raw.Error())
	}
	if len(oe.other) != 0 {
		builder.WriteString("(")
		builder.WriteString(oe.other)
		builder.WriteString(")")
	}
	oe.cache = builder.String()
	return oe.cache
}
