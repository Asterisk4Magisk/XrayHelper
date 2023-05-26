//go:build !linux

package common

import "XrayHelper/main/errors"

// SetUidGid not implement
func (this *external) SetUidGid(uid string, gid string) error {
	return errors.New("system not support SetUidGid").WithPrefix("external_other").WithPathObj(*this)
}
