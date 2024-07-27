package cgroup

import (
	"XrayHelper/main/builds"
	e "XrayHelper/main/errors"
	"XrayHelper/main/log"
	"bufio"
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
			return "", e.New("bad mount entry ", text).WithPrefix(tagCgroup)
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
		memLimit, _ := strconv.ParseFloat(builds.Config.XrayHelper.MemLimit, 64)
		mp, err := v1MountPoint()
		if err != nil {
			return err
		}
		mountPoint = mp
		// create cpu limit
		if cpuLimit != 100.0 {
			if err := os.MkdirAll(filepath.Join(mountPoint, "cpuctl", name), 0o755); err != nil {
				return e.New("cannot create cpuctl cgroup, ", err).WithPrefix(tagCgroup)
			}
			if err := os.WriteFile(
				filepath.Join(mountPoint, "cpuctl", name, "cpu.uclamp.max"),
				[]byte(strconv.FormatFloat(cpuLimit, 'f', 2, 64)),
				os.FileMode(0),
			); err != nil {
				log.HandleDebug("kernel not support uclamp, skip cpu.uclamp.max")
			}
			if err := os.WriteFile(
				filepath.Join(mountPoint, "cpuctl", name, "cpu.shares"),
				[]byte(strconv.FormatInt(int64(cpuLimit*0.01*1024), 10)),
				os.FileMode(0),
			); err != nil {
				return e.New("cannot apply cpuctl cgroup, ", err).WithPrefix(tagCgroup)
			}
		}
		// create memory limit
		if memLimit > 0 {
			if err := os.MkdirAll(filepath.Join(mountPoint, "memcg", name), 0o755); err != nil {
				return e.New("cannot create memcg cgroup, ", err).WithPrefix(tagCgroup)
			}
			if err := os.WriteFile(
				filepath.Join(mountPoint, "memcg", name, "memory.limit_in_bytes"),
				// convert limit to bytes
				[]byte(strconv.FormatInt(int64(memLimit*1024*1024), 10)),
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
