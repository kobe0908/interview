package err_group

import (
	"context"
	"sync"
)

type token struct {
}

type MyErrGroup struct {
	wait   sync.WaitGroup
	sem    chan token
	err    error
	once   sync.Once
	cancel func()
}

func WithContext(ctx context.Context) (*MyErrGroup, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &MyErrGroup{
		cancel: cancel,
	}, ctx
}

func (m *MyErrGroup) done() {
	if m.sem != nil {
		<-m.sem
	}
	m.wait.Done()
}

func (m *MyErrGroup) Go(f func() error) {
	if m.sem != nil {
		m.sem <- token{}
	}
	m.wait.Add(1)
	go func() {
		defer m.done()
		if err := f(); err != nil {
			m.once.Do(func() {
				m.err = err
				if m.cancel != nil {
					m.cancel()
				}
			})
		}
	}()
}

func (m *MyErrGroup) Wait() error {
	m.wait.Wait()
	if m.cancel != nil {
		m.cancel()
	}
	return m.err
}

func (m *MyErrGroup) TryGo(f func() error) bool {
	if m.sem != nil {
		select {
		case m.sem <- token{}:
		default:
			return false
		}
	}
	m.wait.Add(1)
	go func() {
		defer m.done()
		if err := f(); err != nil {
			m.once.Do(func() {
				m.err = err
				if m.cancel != nil {
					m.cancel()
				}
			})
		}
	}()
	return true
}

func (m *MyErrGroup) SetLimit(n int) {
	if n <= 0 {
		return
	}
	m.sem = make(chan token, n)
}
