package workermanager

type Stats struct {
	Workers           int
	QueueLength       int
	MessagesProcessed int
	MessagesTotal     int
}
