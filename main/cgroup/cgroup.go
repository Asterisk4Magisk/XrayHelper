package cgroup

import (
	"XrayHelper/main/builds"
	e "XrayHelper/main/errors"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	tagCgroup = "cgroup"
	name      = "proxy"
)

var mountPoint string

// v1MountPoint returns the mount point where the cgroup
// mountpoints are mounted in a single hierarchy
func v1MountPoint() (string, error) {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return "", err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var (
			text      = scanner.Text()
			fields    = strings.Split(text, " ")
			numFields = len(fields)
		)
		if numFields < 10 {
			return "", fmt.Errorf("mountinfo: bad entry %q", text)
		}
		if fields[numFields-3] == "cgroup" {
			return filepath.Dir(fields[4]), nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", e.New("cgroup v1 mount point not found").WithPrefix(tagCgroup)
}

// LimitProcess use cgroup v1 to limit process resource
func LimitProcess(pid int) error {
	if mountPoint == "" {
		cpuLimit, _ := strconv.ParseFloat(builds.Config.XrayHelper.CPULimit, 64)
		memLimit, _ := strconv.ParseInt(builds.Config.XrayHelper.MemLimit, 10, 64)
		mp, err := v1MountPoint()
		if err != nil {
			return err
		}
		mountPoint = mp
		// create cpu limit
		if cpuLimit != 1.0 {
			uclampMax := int64(float64(1024) * cpuLimit)
			if err := os.MkdirAll(filepath.Join(mountPoint, "cpuctl", name), 0o755); err != nil {
				return e.New("cannot create cpuctl cgroup, ", err).WithPrefix(tagCgroup)
			}
			if err := os.WriteFile(
				filepath.Join(mountPoint, "cpuctl", name, "cpu.uclamp.max"),
				[]byte(strconv.FormatInt(uclampMax, 10)),
				os.FileMode(0),
			); err != nil {
				return e.New("cannot apply cpuctl cgroup, ", err).WithPrefix(tagCgroup)
			}
		}
		// create memory limit
		if memLimit > 0 {
			// convert limit to bytes
			memLimitBytes := memLimit * 1024 * 1024
			if err := os.MkdirAll(filepath.Join(mountPoint, "memcg", name), 0o755); err != nil {
				return e.New("cannot create memcg cgroup, ", err).WithPrefix(tagCgroup)
			}
			if err := os.WriteFile(
				filepath.Join(mountPoint, "memcg", name, "memory.limit_in_bytes"),
				[]byte(strconv.FormatInt(memLimitBytes, 10)),
				os.FileMode(0),
			); err != nil {
				return e.New("cannot apply memcg cgroup, ", err).WithPrefix(tagCgroup)
			}
		}

	}
	// apply cpu limit
	f, err := os.OpenFile(filepath.Join(mountPoint, "cpuctl", name, "cgroup.procs"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0))
	if err != nil {
		return e.New("cannot open cpuctl cgroup procs, ", err).WithPrefix(tagCgroup)
	}
	defer f.Close()
	if _, err := f.WriteString(strconv.FormatInt(int64(pid), 10) + "\n"); err != nil {
		return e.New("cannot add process to cpuctl cgroup, ", err).WithPrefix(tagCgroup)
	}
	// apply memory limit
	f2, err := os.OpenFile(filepath.Join(mountPoint, "memcg", name, "cgroup.procs"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0))
	if err != nil {
		return e.New("cannot open memcg cgroup procs, ", err).WithPrefix(tagCgroup)
	}
	defer f2.Close()
	if _, err := f2.WriteString(strconv.FormatInt(int64(pid), 10) + "\n"); err != nil {
		return e.New("cannot add process to memcg cgroup, ", err).WithPrefix(tagCgroup)
	}
	return nil
}
