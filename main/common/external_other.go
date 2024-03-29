//go:build !linux

package common

import e "XrayHelper/main/errors"

// SetUidGid not implement
func (this *external) SetUidGid(uid string, gid string) error {
	return e.New("system not support SetUidGid").WithPrefix(tagExternal).WithPathObj(*this)
}
