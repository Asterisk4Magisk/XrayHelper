package utils

import (
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"time"
)

type External interface {
	Err() error
	Cmd() *exec.Cmd
	Run()
	RunWithTimeout(timeout time.Duration)
	Start()
	StartWithTimeout(timeout time.Duration)
}

type external struct {
	name   string
	arg    []string
	cmd    *exec.Cmd
	stdout io.Writer
	stderr io.Writer
	err    error
}

func NewExternal(stdout io.Writer, stderr io.Writer, name string, arg ...string) External {
	return &external{name: name, arg: arg, stdout: stdout, stderr: stderr}
}

func (this *external) Run() {
	if this.stdout == nil {
		this.stdout = os.Stdout
	}
	if this.stderr == nil {
		this.stderr = os.Stderr
	}
	this.cmd = exec.Command(this.name, this.arg...)
	this.cmd.Stdout = this.stdout
	this.cmd.Stderr = this.stderr
	this.err = this.cmd.Run()
}

func (this *external) RunWithTimeout(timeout time.Duration) {
	if this.stdout == nil {
		this.stdout = os.Stdout
	}
	if this.stderr == nil {
		this.stderr = os.Stderr
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	this.cmd = exec.CommandContext(ctx, this.name, this.arg...)
	this.cmd.Stdout = this.stdout
	this.cmd.Stderr = this.stderr
	this.err = this.cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		this.err = errors.New("command timed out")
	}
}

func (this *external) Start() {
	if this.stdout == nil {
		this.stdout = os.Stdout
	}
	if this.stderr == nil {
		this.stderr = os.Stderr
	}
	this.cmd = exec.Command(this.name, this.arg...)

	this.cmd.Stdout = this.stdout
	this.cmd.Stderr = this.stderr
	this.err = this.cmd.Start()
}

func (this *external) StartWithTimeout(timeout time.Duration) {
	if this.stdout == nil {
		this.stdout = os.Stdout
	}
	if this.stderr == nil {
		this.stderr = os.Stderr
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	this.cmd = exec.CommandContext(ctx, this.name, this.arg...)
	this.cmd.Stdout = this.stdout
	this.cmd.Stderr = this.stderr
	this.err = this.cmd.Start()
	if ctx.Err() == context.DeadlineExceeded {
		this.err = errors.New("command timed out")
	}
}

func (this *external) Cmd() *exec.Cmd {
	return this.cmd
}

func (this *external) Err() error {
	return this.err
}
