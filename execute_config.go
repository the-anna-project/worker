package worker

import (
	"github.com/the-anna-project/context"
)

// ExecuteConfig represents the configuration used to execute a new worker pool.
type ExecuteConfig struct {
	// Settings.

	// Actions represents the function executed by worker goroutines. The given
	// context is the context provided to the call to Service.Execute.
	Actions []func(ctx context.Context) error
	// CancelOnError defines whether to signal cancelation of worker goroutines in
	// case one worker of the pool throws an error.
	CancelOnError bool
	// Errors is the channel used to put all occured errors into, if any. When
	// this is provided manually the caller have to make sure the errors channel
	// is buffered according to the number of provided actions. TODO check if that is correct
	Errors chan error
	// NumWorkers represents the number of workers to be registered to run
	// concurrently within the pool.
	NumWorkers int
}
