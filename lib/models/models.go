package models

import (
	"time"
	"github.com/coffee-machine/lib/errors"
)

// Future is used to return response for async execution
type Future struct {
	err chan error
	res chan interface{}
}

// NewFuture returns a new instance of empty future object
func NewFuture() *Future {
	return &Future{
		err: make(chan error, 1),
		res: make(chan interface{}, 1),
	}
}

// Error returns the error from the Future object
func (f *Future) Error() error {
	return <-f.err
}

// Result returns the output of the future object
func (f *Future) Result() interface{} {
	return <-f.res
}

// NotifyError notifies the error signal to the future object
func (f *Future) NotifyError(err error) {
	f.err <- err
	close(f.err)
}

// NotifyResult makes the result available to future objedt
func (f *Future) NotifyResult(res interface{}) {
	f.res <- res
	close(f.res)
}

// Task represent a unit of work that can be done independently
type Task struct {
	// CreationTime records the time of creation for the task
	CreationTime int64
	// ttl is time to live for the task, TTL 0 implies infinite
	ttl int64
	// Age records the retry age of the task, i.e the number of time task has been
	// through execution
	Age int
	// f is the unit of task that has to be run
	f func() error
}

// TaskOption represent an option type to override default task attributes
type TaskOption func(t *Task)

// WithTTL sets the TTL for the task
func WithTTL(ttl int64) TaskOption {
	return func(t *Task) {
		t.ttl = ttl
	}
}

// Run executes the business logic for the task, it returns an error if the TTL for the task has expired
func (t *Task) Run() error {
	if t.expiredTTL() {
		return errors.New(ErrTaskExpired, fmt.Errord("task expired"))s
	}

	return t.f()
}

func (t *Task) expiredTTL() bool {
	return t.ttl != 0 && time.Now().Unix() < (t.CreationTime+t.ttl)
}

// NewTask returns a new task with default or provided options
func NewTask(f func() error, options ...TaskOption) *Task {
	task := &Task{
		f: f,
	}
	for _, option := range options {
		option(task)
	}
	return task
}
