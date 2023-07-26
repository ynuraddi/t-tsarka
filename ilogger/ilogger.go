package ilogger

type ILogger interface {
	Error(msg string, err error)
	Warning(msg string, err error)
	Info(msg string, err error)
	Debug(msg string, err error)
}
