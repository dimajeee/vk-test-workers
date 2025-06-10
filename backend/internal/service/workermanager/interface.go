package workermanager

import "context"

type WorkerManager interface {
	AddWorkerWithContext(ctx context.Context)
	RemoveWorker() bool
	Send(msg string) bool
	Stats() (workerCount int, queueLength int)
	StopAll()
	Wait()
	GetStats() Stats
}
