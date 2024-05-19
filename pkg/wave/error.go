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

	builder := strings.Builder{}
	builder.Grow(len("wave: "))
	builder.WriteString("wave: ")
	if len(oe.op) != 0 {
		builder.Grow(len(oe.op))
		builder.WriteString(oe.op)
	}

	if len(oe.port) != 0 {
		builder.Grow(len(oe.port) + 2)
		builder.WriteString("[")
		builder.WriteString(oe.port)
		builder.WriteString("]")
	}
	if oe.raw != nil {
		e := oe.raw.Error()
		builder.Grow(len(e) + 2)
		builder.WriteString(": ")
		builder.WriteString(e)
	}
	if len(oe.other) != 0 {
		builder.Grow(2 + len(oe.other))
		builder.WriteString("(")
		builder.WriteString(oe.other)
		builder.WriteString(")")
	}
	oe.cache = builder.String()
	return oe.cache
}
