//go:build js && wasm

package sdk

import "syscall"

func sysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}
