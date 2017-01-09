package worker

import (
	"testing"

	"github.com/the-anna-project/context"
)

func Test_Execute_OneAction_OneWorker(t *testing.T) {
	newService := testNewService(t)

	actions := []func(ctx context.Context) error{
		func(ctx context.Context) error {
			return nil
		},
	}

	ctx, err := context.New(context.DefaultConfig())
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	config := newService.ExecuteConfig()
	config.Actions = actions
	config.NumWorkers = 1
	err = newService.Execute(ctx, config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
}

func testNewService(t *testing.T) Service {
	newService, err := NewService(DefaultServiceConfig())
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	return newService
}
