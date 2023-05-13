package utils

import (
	"context"
	"errors"
	"io"
	"os/exec"
	"time"
)

type External interface {
	Err() error
	SetUidGid(uid uint32, gid uint32) error
	Run()
	Start()
	Wait() error
	Kill() error
}

type external struct {
	timeout time.Duration
	ctx     context.Context
	cancel  context.CancelFunc
	cmd     *exec.Cmd
	err     error
}

func NewExternal(timeout time.Duration, out io.Writer, err io.Writer, name string, arg ...string) External {
	var e = external{timeout: timeout}
	if timeout > 0 {
		e.ctx, e.cancel = context.WithTimeout(context.Background(), timeout)
		e.cmd = exec.CommandContext(e.ctx, name, arg...)
		e.cmd.Stdout = out
		e.cmd.Stderr = err
	} else {
		e.cmd = exec.Command(name, arg...)
		e.cmd.Stdout = out
		e.cmd.Stderr = err
	}
	return &e
}

func (this *external) Run() {
	if this.timeout > 0 {
		defer this.cancel()
	}
	this.err = this.cmd.Run()
	if this.timeout > 0 {
		if this.ctx.Err() == context.DeadlineExceeded {
			this.err = errors.New("command timed out")
		}
	}
}

func (this *external) Start() {
	this.err = this.cmd.Start()
}

func (this *external) Err() error {
	return this.err
}

func (this *external) Wait() error {
	if this.timeout > 0 {
		defer this.cancel()
	}
	err := this.cmd.Wait()
	if this.timeout > 0 {
		if this.ctx.Err() == context.DeadlineExceeded {
			this.err = errors.New("command timed out")
		}
	}
	return err
}

func (this *external) Kill() error {
	return this.cmd.Process.Kill()
}
