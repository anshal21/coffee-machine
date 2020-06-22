package models

import (
	"fmt"
	"time"

	"github.com/anshal21/coffee-machine/lib"
	"github.com/anshal21/coffee-machine/lib/errors"
)

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
		return errors.New(lib.ErrTaskExpired, fmt.Errorf("task expired"))
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
