package utils

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"time"
)

type External interface {
	Err() error
	Cmd() *exec.Cmd
	Stdout() *bytes.Buffer
	Stderr() *bytes.Buffer
	Run()
	RunWithTimeout(timeout time.Duration)
	Start()
	StartWithTimeout(timeout time.Duration)
}

type external struct {
	name   string
	arg    []string
	cmd    *exec.Cmd
	stdout *bytes.Buffer
	stderr *bytes.Buffer
	err    error
}

func NewExternal(name string, arg ...string) External {
	return &external{name: name, arg: arg}
}

func (this *external) Run() {
	this.cmd = exec.Command(this.name, this.arg...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	this.stdout = &stdout
	this.stderr = &stderr
	this.cmd.Stdout = &stdout
	this.cmd.Stderr = &stderr
	this.err = this.cmd.Run()
}

func (this *external) RunWithTimeout(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	this.cmd = exec.CommandContext(ctx, this.name, this.arg...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	this.stdout = &stdout
	this.stderr = &stderr
	this.cmd.Stdout = &stdout
	this.cmd.Stderr = &stderr
	this.err = this.cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		this.err = errors.New("command timed out")
	}
}

func (this *external) Start() {
	this.cmd = exec.Command(this.name, this.arg...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	this.stdout = &stdout
	this.stderr = &stderr
	this.cmd.Stdout = &stdout
	this.cmd.Stderr = &stderr
	this.err = this.cmd.Start()
}

func (this *external) StartWithTimeout(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	this.cmd = exec.CommandContext(ctx, this.name, this.arg...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	this.stdout = &stdout
	this.stderr = &stderr
	this.cmd.Stdout = &stdout
	this.cmd.Stderr = &stderr
	this.err = this.cmd.Start()
	if ctx.Err() == context.DeadlineExceeded {
		this.err = errors.New("command timed out")
	}
}

func (this *external) Cmd() *exec.Cmd {
	return this.cmd
}

func (this *external) Stdout() *bytes.Buffer {
	return this.stdout
}

func (this *external) Stderr() *bytes.Buffer {
	return this.stderr
}

func (this *external) Err() error {
	return this.err
}
