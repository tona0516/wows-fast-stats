package repository

type ReportInterface interface {
	Send(message string, content error)
}
