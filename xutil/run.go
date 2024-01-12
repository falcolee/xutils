package xutil

import (
	"fmt"

	"github.com/falcolee/xutils/xerror"
)

// SafeRun 能处理 panic 的函数容器
func SafeRun(f func()) (backErr error) {
	defer func() {
		if err := recover(); err != nil {
			msg := xerror.GetCallerInfo(15)
			backErr = fmt.Errorf("panic recover with error %v\n stack is :\n %s", err, msg)
		}
	}()

	f()

	return backErr
}
