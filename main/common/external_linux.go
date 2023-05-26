//go:build linux

package common

import (
	"strconv"
	"syscall"
)

// SetUidGid implement in linux
func (this *external) SetUidGid(uid string, gid string) error {
	uidInt, _ := strconv.Atoi(uid)
	gidInt, _ := strconv.Atoi(gid)
	this.cmd.SysProcAttr = &syscall.SysProcAttr{}
	this.cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uidInt), Gid: uint32(gidInt)}
	return nil
}
