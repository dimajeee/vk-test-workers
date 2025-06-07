package workermanager

type WorkerManager interface {
	AddWorker()
	RemoveWorker() bool
	Send(msg string) bool
	Stats() (workerCount int, queueLength int)
	StopAll()
	Wait()
	GetStats() Stats
}
