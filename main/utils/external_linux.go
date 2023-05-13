//go:build linux

package utils

import "syscall"

// SetUidGid implement in linux
func (this *external) SetUidGid(uid uint32, gid uint32) error {
	this.cmd.SysProcAttr = &syscall.SysProcAttr{}
	this.cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uid, Gid: gid}
	return nil
}
