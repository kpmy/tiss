//import github.com/kpmy/ypk/tc
package tc

import (
	"errors"
	"fmt"
	"sync"
)

type Continue interface {
	Catch(error, func(error)) Continue
	Finally(func()) Continue
	Do(...interface{}) interface{}
}

type catch struct {
	err error
	fn  func(error)
}

type tc struct {
	sync.Once
	fn  func(...interface{}) interface{}
	par []interface{}
	e   []catch
	fin func()
}

func Throw(e ...interface{}) {
	panic(errors.New(fmt.Sprint(e...)))
}

func (t *tc) Catch(e error, fn func(error)) Continue {
	t.e = append(t.e, catch{e, fn})
	return t
}

func (t *tc) Do(par ...interface{}) (ret interface{}) {
	t.Once.Do(func() {
		defer func() {
			if _x := recover(); _x != nil {
				switch x := _x.(type) {
				case error:
					var next func(error)
					for _, c := range t.e {
						if c.err == x {
							next = c.fn
							break
						}
						if c.err == nil {
							next = c.fn
						}
					}
					if next != nil {
						next(x)
					}
					if t.fin != nil {
						t.fin()
					}
					if next == nil {
						panic(x)
					}
				default:
					var next func(error)
					for _, c := range t.e {
						if c.err == nil {
							next = c.fn
							break
						}
					}
					if next != nil {
						err := errors.New(fmt.Sprint(_x))
						next(err)
					} else {
						panic(_x)
					}
				}
			}
		}()
		t.par = append(t.par, par...)
		ret = t.fn(t.par...)
		if t.fin != nil {
			t.fin()
		}
	})
	return
}

func (t *tc) Finally(fin func()) Continue {
	t.fin = fin
	return t
}

func Try(fn func(...interface{}) interface{}, par ...interface{}) Continue {
	ret := &tc{}
	ret.fn = fn
	ret.par = par
	return ret
}

func Do(fn func()) Continue {
	return Try(func(...interface{}) interface{} {
		fn()
		return nil
	})
}
