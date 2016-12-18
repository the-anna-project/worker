package worker

// Service represents a service to process work concurrently.
type Service interface {
	// Execute processes all configured actions concurrently. The call to Execute
	// blocks until all goroutines within the worker pool have finished their
	// work.
	Execute(config ExecuteConfig) error
	// ExecuteConfig provides a default configuration for Execute.
	ExecuteConfig() ExecuteConfig
}
