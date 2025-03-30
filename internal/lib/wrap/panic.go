package wrap

import (
	"fmt"
	"runtime/debug"
)

type PanicError struct {
	err   error
	stack string
}

func (e PanicError) Error() string {
	return e.err.Error()
}

func (e PanicError) Unwrap() error {
	return e.err
}

func Panic(fn func() error) (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = PanicError{
				err:   fmt.Errorf("panic: %v", v),
				stack: string(debug.Stack()),
			}
		}
	}()
	return fn()
}
