package dao

import "strings"

type OpErr struct {
	op         string
	database   string
	err        error
	suggestion string
}

func (o *OpErr) Error() string {
	if o == nil {
		return "<nil>"
	}

	es := ""
	if o.err != nil {
		es = o.err.Error()
	}

	return strings.Join([]string{"database: ", o.op, o.database, es, o.suggestion}, " ")
}

func (o *OpErr) Unwrap() error {
	return o.err
}
