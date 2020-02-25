// Package event for internal and protocol events
package event

import (
	"sync/atomic"
)

// Internal event
type Internal struct {
	complete uint32
	ch       chan error
	err      error
	data     interface{}
}

// Result blocks until event is complete
func (ev *Internal) Result() error {
	for atomic.CompareAndSwapUint32(&ev.complete, 1, 1) {
		if atomic.CompareAndSwapUint32(&ev.complete, 2, 2) {
			return ev.err
		}
	}
	<-ev.ch

	return ev.err
}

// Done indicates event completed processing
func (ev *Internal) Done(err error) {
	for atomic.CompareAndSwapUint32(&ev.complete, 0, 1) {
		ev.err = err
		close(ev.ch)
		atomic.CompareAndSwapUint32(&ev.complete, 1, 2)
	}
}

// Data returns event data
func (ev *Internal) Data() interface{} {
	return ev.data
}

// Protocol event structure
type Protocol struct {
	ID   string      `json:"id"`
	Kind string      `json:"kind"`
	Data interface{} `json:"data"`
}

// NewInternal returns new internal event
func NewInternal(value interface{}) *Internal {
	return &Internal{
		complete: 0,
		ch:       make(chan error),
		err:      nil,
		data:     value,
	}
}

// SendInternal sends new internal event
func SendInternal(ch chan<- *Internal, value interface{}) error {
	ev := NewInternal(value)
	ch <- ev
	return ev.Result()
}
