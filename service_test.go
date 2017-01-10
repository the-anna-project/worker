package worker

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/the-anna-project/context"
)

func Test_Execute_OneAction_OneWorker(t *testing.T) {
	newService := testNewService(t)

	var f string
	var l []string
	var m sync.Mutex
	var o sync.Once

	actions := []func(ctx context.Context) error{
		func(ctx context.Context) error {
			k := "one"

			m.Lock()
			defer m.Unlock()
			l = append(l, k)

			o.Do(func() {
				f = k
			})

			return nil
		},
	}

	ctx, err := context.New(context.DefaultConfig())
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	config := newService.ExecuteConfig()
	config.Actions = actions
	config.CancelOnError = true
	config.NumWorkers = 1
	err = newService.Execute(ctx, config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	if len(l) != 1 {
		t.Fatal("expected", 1, "got", len(l))
	}
	if l[0] != "one" {
		t.Fatal("expected", "one", "got", l[0])
	}
}

func Test_Execute_OneActions_ThreeWorker(t *testing.T) {
	newService := testNewService(t)

	var f string
	var l []string
	var m sync.Mutex
	var o sync.Once

	actions := []func(ctx context.Context) error{
		func(ctx context.Context) error {
			k := "one"

			m.Lock()
			defer m.Unlock()
			l = append(l, k)

			o.Do(func() {
				f = k
			})

			return nil
		},
	}

	ctx, err := context.New(context.DefaultConfig())
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	config := newService.ExecuteConfig()
	config.Actions = actions
	config.CancelOnError = true
	config.NumWorkers = 3
	err = newService.Execute(ctx, config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	if len(l) != 1 {
		t.Fatal("expected", 1, "got", len(l))
	}
	if l[0] != "one" {
		t.Fatal("expected", "one", "got", l[0])
	}
}

func Test_Execute_ThreeActions_OneWorker(t *testing.T) {
	newService := testNewService(t)

	var f string
	var l []string
	var m sync.Mutex
	var o sync.Once

	actions := []func(ctx context.Context) error{
		func(ctx context.Context) error {
			time.Sleep(10 * time.Millisecond)

			k := "one"

			m.Lock()
			defer m.Unlock()
			l = append(l, k)

			o.Do(func() {
				f = k
			})

			return nil
		},
		func(ctx context.Context) error {
			time.Sleep(5 * time.Millisecond)

			k := "two"

			m.Lock()
			defer m.Unlock()
			l = append(l, k)

			o.Do(func() {
				f = k
			})

			return nil
		},
		func(ctx context.Context) error {
			time.Sleep(1 * time.Millisecond)

			k := "three"

			m.Lock()
			defer m.Unlock()
			l = append(l, k)

			o.Do(func() {
				f = k
			})

			return nil
		},
	}

	ctx, err := context.New(context.DefaultConfig())
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	config := newService.ExecuteConfig()
	config.Actions = actions
	config.CancelOnError = true
	config.NumWorkers = 1
	err = newService.Execute(ctx, config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// When one worker executes all actions, we expect all action identifiers of
	// the test to be present in series.
	if len(l) != 3 {
		t.Fatal("expected", 3, "got", len(l))
	}
	if l[0] != "one" {
		t.Fatal("expected", "one", "got", l[0])
	}
	if l[1] != "two" {
		t.Fatal("expected", "two", "got", l[1])
	}
	if l[2] != "three" {
		t.Fatal("expected", "three", "got", l[2])
	}

	// When one worker executes all actions, we expect the list of actions to be
	// executed sequentially. That means that the first action is always called
	// first, because there is only one worker wich picks the first from the list.
	if f != "one" {
		t.Fatal("expected", "three", "got", f)
	}
}

func Test_Execute_ThreeActions_ThreeWorker(t *testing.T) {
	newService := testNewService(t)

	var f string
	var l []string
	var m sync.Mutex
	var o sync.Once

	actions := []func(ctx context.Context) error{
		func(ctx context.Context) error {
			time.Sleep(10 * time.Millisecond)

			k := "one"

			m.Lock()
			defer m.Unlock()
			l = append(l, k)

			o.Do(func() {
				f = k
			})

			return nil
		},
		func(ctx context.Context) error {
			time.Sleep(5 * time.Millisecond)

			k := "two"

			m.Lock()
			defer m.Unlock()
			l = append(l, k)

			o.Do(func() {
				f = k
			})

			return nil
		},
		func(ctx context.Context) error {
			time.Sleep(1 * time.Millisecond)

			k := "three"

			m.Lock()
			defer m.Unlock()
			l = append(l, k)

			o.Do(func() {
				f = k
			})

			return nil
		},
	}

	ctx, err := context.New(context.DefaultConfig())
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	config := newService.ExecuteConfig()
	config.Actions = actions
	config.CancelOnError = true
	config.NumWorkers = 3
	err = newService.Execute(ctx, config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// When three workers execute all actions, we expect all action identifiers of
	// the test to be present.
	if len(l) != 3 {
		t.Fatal("expected", 3, "got", len(l))
	}
	if !testContainsString(l, "one") {
		t.Fatal("expected", true, "got", false)
	}
	if !testContainsString(l, "two") {
		t.Fatal("expected", true, "got", false)
	}
	if !testContainsString(l, "three") {
		t.Fatal("expected", true, "got", false)
	}

	// When three workers execute all actions, we expect the list of actions to be
	// executed concurrently. That means that the fastest action is always
	// finished first, because there is one worker for each action wich makes all
	// actions of the list being executed simultaneously.
	if f != "three" {
		t.Fatal("expected", "three", "got", f)
	}
}

func Test_Execute_OneAction_OneWorker_Error(t *testing.T) {
	newService := testNewService(t)

	var f error
	var l []error
	var m sync.Mutex
	var o sync.Once

	actions := []func(ctx context.Context) error{
		func(ctx context.Context) error {
			e := fmt.Errorf("one")

			m.Lock()
			defer m.Unlock()
			l = append(l, e)

			o.Do(func() {
				f = e
			})

			return e
		},
	}

	ctx, err := context.New(context.DefaultConfig())
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	config := newService.ExecuteConfig()
	config.Actions = actions
	config.CancelOnError = true
	config.NumWorkers = 1
	err = newService.Execute(ctx, config)
	if err.Error() != "one" {
		t.Fatal("expected", "one", "got", err.Error())
	}

	if len(l) != 1 {
		t.Fatal("expected", 1, "got", len(l))
	}
	if !testContainsError(l, "one") {
		t.Fatal("expected", true, "got", false)
	}
}

func Test_Execute_ThreeActions_ThreeWorker_Error(t *testing.T) {
	newService := testNewService(t)

	var f error
	var l []error
	var m sync.Mutex
	var o sync.Once

	actions := []func(ctx context.Context) error{
		func(ctx context.Context) error {
			var e error

			select {
			case <-ctx.Done():
			case <-time.After(10 * time.Millisecond):
				e = fmt.Errorf("one")

				m.Lock()
				defer m.Unlock()
				l = append(l, e)

				o.Do(func() {
					f = e
				})
			}

			return e
		},
		func(ctx context.Context) error {
			var e error

			select {
			case <-ctx.Done():
			case <-time.After(5 * time.Millisecond):
				e = fmt.Errorf("two")

				m.Lock()
				defer m.Unlock()
				l = append(l, e)

				o.Do(func() {
					f = e
				})
			}

			return e
		},
		func(ctx context.Context) error {
			var e error

			select {
			case <-ctx.Done():
			case <-time.After(1 * time.Millisecond):
				e = fmt.Errorf("three")

				m.Lock()
				defer m.Unlock()
				l = append(l, e)

				o.Do(func() {
					f = e
				})
			}

			return e
		},
	}

	ctx, err := context.New(context.DefaultConfig())
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	config := newService.ExecuteConfig()
	config.Actions = actions
	config.CancelOnError = true
	config.NumWorkers = 3
	err = newService.Execute(ctx, config)
	if err.Error() != "three" {
		t.Fatal("expected", "three", "got", err.Error())
	}

	// We cancel on the first error, which is produced by the third action.
	if len(l) != 1 {
		t.Fatal("expected", 1, "got", len(l))
	}
	if !testContainsError(l, "three") {
		t.Fatal("expected", true, "got", false)
	}

	if f.Error() != "three" {
		t.Fatal("expected", "three", "got", f.Error())
	}
}

func Test_Execute_ThreeActions_ThreeWorker_Error_NoCancelOnError(t *testing.T) {
	newService := testNewService(t)

	var f error
	var l []error
	var m sync.Mutex
	var o sync.Once

	actions := []func(ctx context.Context) error{
		func(ctx context.Context) error {
			var e error

			select {
			case <-ctx.Done():
			case <-time.After(10 * time.Millisecond):
				e = fmt.Errorf("one")

				m.Lock()
				defer m.Unlock()
				l = append(l, e)

				o.Do(func() {
					f = e
				})
			}

			return e
		},
		func(ctx context.Context) error {
			var e error

			select {
			case <-ctx.Done():
			case <-time.After(5 * time.Millisecond):
				e = fmt.Errorf("two")

				m.Lock()
				defer m.Unlock()
				l = append(l, e)

				o.Do(func() {
					f = e
				})
			}

			return e
		},
		func(ctx context.Context) error {
			var e error

			select {
			case <-ctx.Done():
			case <-time.After(1 * time.Millisecond):
				e = fmt.Errorf("three")

				m.Lock()
				defer m.Unlock()
				l = append(l, e)

				o.Do(func() {
					f = e
				})
			}

			return e
		},
	}

	ctx, err := context.New(context.DefaultConfig())
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	config := newService.ExecuteConfig()
	config.Actions = actions
	config.CancelOnError = false
	config.NumWorkers = 3
	err = newService.Execute(ctx, config)
	if err.Error() != "three" {
		t.Fatal("expected", "three", "got", err.Error())
	}

	if len(l) != 3 {
		t.Fatal("expected", 3, "got", len(l))
	}
	if !testContainsError(l, "one") {
		t.Fatal("expected", true, "got", false)
	}
	if !testContainsError(l, "two") {
		t.Fatal("expected", true, "got", false)
	}
	if !testContainsError(l, "three") {
		t.Fatal("expected", true, "got", false)
	}

	if f.Error() != "three" {
		t.Fatal("expected", "three", "got", f.Error())
	}
}

func testContainsError(l []error, i string) bool {
	for _, k := range l {
		if k.Error() == i {
			return true
		}
	}

	return false
}

func testContainsString(l []string, i string) bool {
	for _, k := range l {
		if k == i {
			return true
		}
	}

	return false
}

func testNewService(t *testing.T) Service {
	newService, err := NewService(DefaultServiceConfig())
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	return newService
}
