package adapter

type ReportInterface interface {
	Send(message string, content error)
}
