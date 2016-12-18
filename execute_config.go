package worker

// ExecuteConfig represents the configuration used to execute a new worker pool.
type ExecuteConfig struct {
	// Settings.

	// Actions represents the function executed by workers.
	Actions []func(canceler <-chan struct{}) error
	// Canceler can be used to end the worker pool's processes pro-actively. The
	// signal received here will be redirected to the canceler provided to the
	// worker functions.
	Canceler chan struct{}
	// CancelOnError defines whether to signal cancelation of worker processes in
	// case one worker of the pool throws an error.
	CancelOnError bool
	// NumWorkers represents the number of workers to be registered to run
	// concurrently within the pool.
	NumWorkers int
}
