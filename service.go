// Package worker implements a service to process work concurrently.
package worker

import (
	"sync"

	"github.com/the-anna-project/context"
)

// ServiceConfig represents the configuration used to create a new worker
// service.
type ServiceConfig struct {
}

// DefaultServiceConfig provides a default configuration to create a new worker
// service by best effort.
func DefaultServiceConfig() ServiceConfig {
	return ServiceConfig{}
}

// NewService creates a new configured worker service.
func NewService(config ServiceConfig) (Service, error) {
	return &service{}, nil
}

type service struct {
}

func (s *service) Execute(ctx context.Context, config ExecuteConfig) error {
	if config.Errors == nil {
		config.Errors = make(chan error, len(config.Actions))
	}

	var wg sync.WaitGroup

	for _, action := range config.Actions {
		go func() {
			for n := 0; n < config.NumWorkers; n++ {
				wg.Add(1)
				go func(action func(ctx context.Context) error) {
					defer wg.Done()

					err := action(ctx)
					if err != nil {
						// We want to capture errors in any case, so we do this at first.
						config.Errors <- err

						if config.CancelOnError {
							// Canceling the context acts as broadcast to all workers that
							// should listen to the context's done channel.
							ctx.Cancel()
						}
					}
				}(action)
			}
		}()
	}

	wg.Wait()

	select {
	case err := <-config.Errors:
		// The errors channel is supposed to hold all errors. Here it also serves as
		// internal state of truth. So we have to read the first error occured.
		// Then, to not modify the errors channel from the client point of view, we
		// put the read error back.
		config.Errors <- err
		return err
	default:
		return nil
	}
}

func (s *service) ExecuteConfig() ExecuteConfig {
	return ExecuteConfig{
		// Settings.

		Actions:       []func(ctx context.Context) error{},
		CancelOnError: true,
		Errors:        nil,
		NumWorkers:    1,
	}
}
