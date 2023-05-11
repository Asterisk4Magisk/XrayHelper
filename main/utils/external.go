package utils

import (
	"bytes"
	"os/exec"
)

type External interface {
	Err() error
	Stdout() string
	Stderr() string
	Run()
}

type external struct {
	cmd    *exec.Cmd
	stdout string
	stderr string
	err    error
}

func NewExternal(cmd *exec.Cmd) External {
	return &external{cmd: cmd}
}

func (this *external) Run() {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	this.cmd.Stdout = &stdout
	this.cmd.Stderr = &stderr
	this.err = this.cmd.Run()
	this.stdout = stdout.String()
	this.stderr = stderr.String()
}

func (this *external) Stdout() string {
	return this.stdout
}

func (this *external) Stderr() string {
	return this.stderr
}

func (this *external) Err() error {
	return this.err
}
