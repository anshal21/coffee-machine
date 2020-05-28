package workerpool

import (
	"sync"

	"github.com/coffee-machine/lib/models"
)

// Request contains the parameters to initialize a worker pool
// Default values are used for missing parameters
// Following are the default values
//	- Size: 1
//	- Buffer: 100
type Request struct {
	Size   int
	Buffer int
}

const (
	_defaultBufferSize  = 1000
	_defaultWorkerCount = 1
)

// WorkerPool is an interface to a workerpool implementation
// TODO: Add support for version of tasks that could return a response as well along with error
type WorkerPool interface {
	// Add adds a task to the worker pool
	// the function returns a future object that holds the response of the execution
	// currently it takes Task as input so only error response is supported
	Add(task *models.Task) *models.Future
	// Start starts the execution of the tasks in the workerpool
	Start()
	// Done notifies the workerpool that there are no more tasks to be added
	Done()
	// Abort notifies workerpool to abort the execution of all remaning tasks in the pool
	Abort()
	// WaitForCompletion is a blocking function that waits till workerpool finishes all the tasks
	WaitForCompletion()
}

type taskWrapper struct {
	task     *models.Task
	response *models.Future
}

// New returns a new instance of WorkerPool with provided specs
func New(request *Request) WorkerPool {
	request = populateDefaults(request)

	return &workerpool{
		workerCount: request.Size,
		abortChan:   make(chan struct{}, request.Size),
		taskChan:    make(chan taskWrapper, request.Buffer),
	}
}

func populateDefaults(request *Request) *Request {
	if request.Size == 0 {
		request.Size = _defaultWorkerCount
	}
	if request.Buffer == 0 {
		request.Buffer = _defaultBufferSize
	}
	return request
}

type workerpool struct {
	wg          sync.WaitGroup
	workerCount int
	abortChan   chan struct{}
	taskChan    chan taskWrapper
}

func (w *workerpool) Start() {
	w.wg.Add(w.workerCount)
	for i := 0; i < w.workerCount; i++ {
		go func() {
			defer w.wg.Done()
			for {
				select {
				case task, ok := <-w.taskChan:
					if !ok {
						return
					}
					err := task.task.Run()
					task.response.NotifyResult(nil)
					task.response.NotifyError(err)
				case <-w.abortChan:
					w.taskChan = nil
					return

				}
				if w.taskChan == nil {
					return
				}
			}
		}()
	}
}

func (w *workerpool) Add(task *models.Task) *models.Future {
	res := models.NewFuture()
	workUnit := taskWrapper{
		task:     task,
		response: res,
	}

	w.taskChan <- workUnit
	return res
}

func (w *workerpool) Done() {
	close(w.taskChan)
}

func (w *workerpool) Abort() {
	for i := 0; i < w.workerCount; i++ {
		w.abortChan <- struct{}{}
	}
}

func (w *workerpool) WaitForCompletion() {
	w.wg.Wait()
