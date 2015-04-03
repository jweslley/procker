// +build darwin dragonfly freebsd netbsd openbsd

package main

import "syscall"

func sysProcAttrs() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setpgid: true,
	}
}
