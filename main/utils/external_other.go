//go:build !linux

package utils

import "errors"

func (this *external) SetUidGid(uid uint32, gid uint32) error {
	return errors.New("SetUidGid: system not support")
}
