//go:build !linux

package utils

import "XrayHelper/main/errors"

// SetUidGid not implement
func (this *external) SetUidGid(uid uint32, gid uint32) error {
	return errors.New("system not support SetUidGid").WithPrefix("external_other").WithPathObj(*this)
}
