// Package worker implements a service to process work concurrently.
package worker

import (
	"sync"
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

func (s *service) Execute(config ExecuteConfig) error {
	var wg sync.WaitGroup
	var once sync.Once

	canceler := make(chan struct{}, 1)

	if config.Canceler != nil {
		go func() {
			<-config.Canceler
			// Receiving a signal from the global canceler will forward the
			// cancelation to all workers. Simply closing the workers canceler wil
			// broadcast the signal to each listener. Here we also make sure we do
			// not close on a closed channel by only closing once.
			once.Do(func() {
				close(canceler)
			})
		}()
	}

	for n := 0; n < config.NumWorkers; n++ {
		go func() {
			for _, action := range config.Actions {
				wg.Add(1)
				go func(action func(canceler <-chan struct{}) error) {
					defer wg.Done()

					err := action(canceler)
					if err != nil {
						// We want to capture errors in any case, so we do this at first.
						config.Errors <- err

						if config.CancelOnError && config.Canceler != nil {
							// Closing the canceler channel acts as broadcast to all workers that
							// should listen to the canceler. Here we also make sure we do not
							// close on a closed channel by only closing once.
							once.Do(func() {
								close(config.Canceler)
							})
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

		Actions:       []func(canceler <-chan struct{}) error{},
		Canceler:      nil,
		CancelOnError: true,
		Errors:        make(chan error, 1),
		NumWorkers:    1,
	}
}
