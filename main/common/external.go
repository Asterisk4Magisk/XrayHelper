package common

import (
	e "XrayHelper/main/errors"
	"context"
	"errors"
	"io"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

type External interface {
	Err() error
	SetUidGid(uid string, gid string)
	AppendEnv(env string)
	Run()
	Start()
	Pid() int
	Wait() error
	Kill() error
}

const tagExternal = "external"

// external implement the interface External, wrapping of exec, easier to use
type external struct {
	timeout time.Duration
	ctx     context.Context
	cancel  context.CancelFunc
	cmd     *exec.Cmd
	err     error
}

// NewExternal returns a new external object with cmd
func NewExternal(timeout time.Duration, out io.Writer, err io.Writer, name string, arg ...string) External {
	var ex = external{timeout: timeout}
	if timeout > 0 {
		ex.ctx, ex.cancel = context.WithTimeout(context.Background(), timeout)
		ex.cmd = exec.CommandContext(ex.ctx, name, arg...)
		ex.cmd.Stdout = out
		ex.cmd.Stderr = err
	} else {
		ex.cmd = exec.Command(name, arg...)
		ex.cmd.Stdout = out
		ex.cmd.Stderr = err
	}
	ex.cmd.SysProcAttr = &syscall.SysProcAttr{}
	ex.cmd.SysProcAttr.Setpgid = true
	return &ex
}

// SetUidGid implement in linux
func (this *external) SetUidGid(uid string, gid string) {
	uidInt, _ := strconv.Atoi(uid)
	gidInt, _ := strconv.Atoi(gid)
	this.cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uidInt), Gid: uint32(gidInt)}
}

// AppendEnv add env variable, eg: JAVA_HOME=/usr/local/java/
func (this *external) AppendEnv(env string) {
	this.cmd.Env = append(this.cmd.Env, env)
}

func (this *external) Run() {
	if this.timeout > 0 {
		defer this.cancel()
	}
	this.err = this.cmd.Run()
	if this.timeout > 0 {
		if errors.Is(this.ctx.Err(), context.DeadlineExceeded) {
			this.err = e.New("command timed out").WithPrefix(tagExternal).WithPathObj(*this)
		}
	}
}

func (this *external) Start() {
	this.err = this.cmd.Start()
}

// Err get the external cmd error
func (this *external) Err() error {
	return this.err
}

// Pid get the external cmd pid
func (this *external) Pid() int {
	return this.cmd.Process.Pid
}

// Wait block thread to wait external cmd complete, external should be started by Start
func (this *external) Wait() error {
	if this.timeout > 0 {
		defer this.cancel()
	}
	this.err = this.cmd.Wait()
	if this.timeout > 0 {
		if errors.Is(this.ctx.Err(), context.DeadlineExceeded) {
			this.err = e.New("command timed out").WithPrefix(tagExternal).WithPathObj(*this)
		}
	}
	return this.err
}

func (this *external) Kill() error {
	return this.cmd.Process.Kill()
}
