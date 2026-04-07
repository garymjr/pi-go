package rpc

import (
	"io"
	"time"
)

type Options struct {
	Executable        string
	Dir               string
	Env               []string
	Provider          string
	Model             string
	NoSession         bool
	SessionDir        string
	ExtraArgs         []string
	ResponseTimeout   time.Duration
	ShutdownTimeout   time.Duration
	StartupDelay      time.Duration
	Stderr            io.Writer
	StderrBufferBytes int
}

func (o Options) withDefaults() Options {
	if o.Executable == "" {
		o.Executable = "pi"
	}
	if o.ResponseTimeout <= 0 {
		o.ResponseTimeout = 30 * time.Second
	}
	if o.ShutdownTimeout <= 0 {
		o.ShutdownTimeout = 2 * time.Second
	}
	if o.StartupDelay <= 0 {
		o.StartupDelay = 100 * time.Millisecond
	}
	if o.StderrBufferBytes <= 0 {
		o.StderrBufferBytes = 64 * 1024
	}
	return o
}
