package tasks

import (
	"context"
	"sync"
)

// Task is a long-running function that honors ctx cancellation.
type Task func(ctx context.Context) error

// Manager supervises a set of tasks & waits for them to stop.
type Manager struct {
	wg    sync.WaitGroup
	tasks []Task
}

// New returns an empty manager.
func New() *Manager { return &Manager{} }

// Add registers a task; call before Start.
func (m *Manager) Add(t Task) { m.tasks = append(m.tasks, t) }

// Start launches every task in its own goroutine.
func (m *Manager) Start(ctx context.Context) {
	for _, t := range m.tasks {
		m.wg.Add(1)
		go func(task Task) {
			defer m.wg.Done()
			_ = task(ctx) // tasks handle their own errors/logs
		}(t)
	}
}

// Wait blocks until all tasks return (Start must be called first).
func (m *Manager) Wait() { m.wg.Wait() }
