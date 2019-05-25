package main

import (
	"github.com/m9rco/exile/kernel/common"
	"github.com/m9rco/exile/kernel/worker"
	"runtime"
	"sync"
)

func main() {
	// Initialize the environment
	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup
	var (
		err error
	)

	// Initialize the container
	if err = common.InitContainer(); err != nil {
		goto ERROR
	}
	// Initialize the serve register
	if err = worker.InitRegister(); err != nil {
		goto ERROR
	}

	// Initialize the serve register
	if err = worker.InitLogSink(); err != nil {
		goto ERROR
	}

	// Initialize the executor
	if err = worker.InitExecutor(); err != nil {
		goto ERROR
	}

	// Initialize the worker job scheduler
	if err = worker.InitScheduler(); err != nil {
		goto ERROR
	}

	// Initialize the worker job manager
	if err = worker.InitJobMgr(); err != nil {
		goto ERROR
	}

	wg.Add(1)
	wg.Wait()
ERROR:
}