package tc

import (
	"errors"
	"fmt"
)

func Assert(cond bool, code int, msg ...interface{}) {
	var e string
	if !cond {
		switch {
		case (code >= 20) && (code < 40):
			e = fmt.Sprint(code, " precondition violated ", fmt.Sprint(msg...))
		case (code >= 40) && (code < 60):
			e = fmt.Sprint(code, " subcondition violated ", fmt.Sprint(msg...))
		case (code >= 60) && (code < 80):
			e = fmt.Sprint(code, " postcondition violated ", fmt.Sprint(msg...))
		default:
			e = fmt.Sprint(code, " ", fmt.Sprint(msg...))
		}
		panic(errors.New(e))
	}
}

func Halt(code int, msg ...interface{}) {
	e := fmt.Sprint(code)
	if len(msg) > 0 {
		e = fmt.Sprint(code, " ", fmt.Sprint(msg...))
	}
	panic(errors.New(e))
}
