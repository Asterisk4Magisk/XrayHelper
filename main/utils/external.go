package utils

import (
	"XrayHelper/main/errors"
	"context"
	"io"
	"net"
	"os/exec"
	"strconv"
	"time"
)

type External interface {
	Err() error
	SetUidGid(uid uint32, gid uint32) error
	Run()
	Start()
	Pid() int
	Wait() error
	Kill() error
}

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
	this.err = errors.New(this.cmd.Run()).WithPrefix("external").WithPathObj(*this.cmd)
	if this.timeout > 0 {
		if this.ctx.Err() == context.DeadlineExceeded {
			this.err = errors.New("command timed out").WithPrefix("external").WithPathObj(*this)
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
	err := errors.New(this.cmd.Wait()).WithPrefix("external").WithPathObj(*this.cmd)
	if this.timeout > 0 {
		if this.ctx.Err() == context.DeadlineExceeded {
			this.err = errors.New("command timed out").WithPrefix("external").WithPathObj(*this)
		}
	}
	return err
}

func (this *external) Kill() error {
	return this.cmd.Process.Kill()
}

// CheckPort check whether the port is listening
func CheckPort(protocol string, host string, port int) bool {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout(protocol, addr, 3*time.Second)
	if err != nil {
		return false
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	return true
}
