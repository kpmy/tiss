package fn

import (
	"fmt"
	"reflect"
)

type Maybe interface {
	String() string
}

type mbn struct{}

func (m *mbn) String() string { return "" }

type mbs struct {
	val interface{}
}

func (m *mbs) String() string { return fmt.Sprint(m.val) }

func MaybeString(ls ...string) (ret Maybe) {
	res := &mbs{val: ""}
	ret = res
	for _, s := range ls {
		if s != "" {
			res.val = fmt.Sprint(res.val, s)
		} else {
			ret = nil
		}
	}
	if ret == nil {
		ret = &mbn{}
	}
	return
}

func IsNil(a interface{}) bool {
	defer func() { recover() }()
	return a == nil || reflect.ValueOf(a).IsNil()
}
