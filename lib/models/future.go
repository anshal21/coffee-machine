package models

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
