package workermanager

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"log/slog"
)

type Manager struct {
	mu                sync.RWMutex
	wg                sync.WaitGroup
	counter           atomic.Uint32
	cancelMap         map[int]context.CancelFunc
	input             chan string
	messagesProcessed atomic.Uint32
	messagesTotal     atomic.Uint32
	closed            bool
}

func New(queueSize int) *Manager {
	return &Manager{
		cancelMap: make(map[int]context.CancelFunc),
		input:     make(chan string, queueSize),
	}
}

func (m *Manager) AddWorker() {
	m.mu.Lock()
	defer m.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	id := int(m.counter.Add(1)) - 1
	m.cancelMap[id] = cancel

	m.wg.Add(1)
	go m.worker(ctx, id)

	slog.Info("Added worker", "id", id)
}

func (m *Manager) RemoveWorker() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.cancelMap) == 0 {
		slog.Warn("No workers to remove")
		return false
	}

	var lastID int
	for id := range m.cancelMap {
		if id > lastID {
			lastID = id
		}
	}

	m.cancelMap[lastID]()
	delete(m.cancelMap, lastID)
	m.counter.Add(1)
	slog.Info("Removed worker", "id", lastID)
	return true
}

func (m *Manager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, cancel := range m.cancelMap {
		cancel()
		slog.Info("Stopped worker via StopAll", "id", id)
	}
	m.cancelMap = make(map[int]context.CancelFunc)
}

func (m *Manager) CloseInput() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.closed {
		close(m.input)
		m.closed = true
	}
}

func (m *Manager) Send(msg string) bool {
	select {
	case m.input <- msg:
		m.messagesTotal.Add(1)
		slog.Debug("Sent message", "msg", msg)
		return true
	default:
		slog.Warn("Input queue full, dropping message", "msg", msg)
		return false
	}
}

func (m *Manager) Stats() (int, int) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.cancelMap), len(m.input)
}

func (m *Manager) Wait() {
	m.wg.Wait()
}

func (m *Manager) worker(ctx context.Context, id int) {
	defer m.wg.Done()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Worker stopped", "id", id)
			return
		case msg, ok := <-m.input:
			if !ok {
				slog.Info("Input channel closed, worker exiting", "id", id)
				return
			}
			m.counter.Add(1)
			slog.Info("Worker received message", "worker_id", id, "msg", msg)
			// Имитация обработки
			time.Sleep(500 * time.Millisecond)
		}
	}
}

type Stats struct {
	Workers           int
	QueueLength       int
	MessagesProcessed int
	MessagesTotal     int
}

func (m *Manager) GetStats() Stats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return Stats{
		Workers:           int(m.counter.Load()),
		QueueLength:       len(m.input),
		MessagesProcessed: int(m.messagesProcessed.Load()),
		MessagesTotal:     int(m.messagesTotal.Load()),
	}
}
