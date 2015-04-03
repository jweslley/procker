// +build windows

package main

import "syscall"

func sysProcAttrs() *syscall.SysProcAttr {
	return nil
}
