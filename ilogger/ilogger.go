package ilogger

type ILogger interface {
	Fatal(msg string, err error)
	Error(msg string, err error)
	Warning(msg string)
	Info(msg string)
	Debug(msg string)
}
