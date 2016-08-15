package trace

import (
	"fmt"
	"io"
)

type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

type nilTracer struct{}

// Trace write a message to a specific output
func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}

// Trace does nothing
func (t *nilTracer) Trace(a ...interface{}) {}

// New creates a tracer that will display messages on a specific output (param)
func New(w io.Writer) *tracer {
	return &tracer{out: w}
}

// Off creates a tracer that won't display anything
func Off() *nilTracer {
	return &nilTracer{}
}
