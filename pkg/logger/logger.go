package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/ynuraddi/t-tsarka/ilogger"
)

var _ ilogger.ILogger = (*Logger)(nil)

type Level int8

const (
	LvlErr Level = iota
	LvlWrn
	LvlInf
	LvlDeb
)

type Logger struct {
	output   io.Writer
	minLevel Level
	mu       sync.Mutex

	sigChan chan os.Signal
}

// need chan osSignal for gracefull shutdown
func NewLogger(out io.Writer, logLevel Level, sigChan chan<- os.Signal) *Logger {
	logger := &Logger{
		output:   out,
		minLevel: Level(logLevel),
		mu:       sync.Mutex{},
	}

	return logger
}

func (l *Logger) print(lvl Level, msg string) {
	if lvl > l.minLevel {
		return
	}

	msg = fmt.Sprintf("%s\t%s", time.Now().Format(time.TimeOnly), msg)

	l.output.Write(append([]byte(msg), '\n'))
}

func (l *Logger) Fatal(msg string, err error) {
	logMsg := fmt.Sprintf("[FATAL]\t%s:\t%v", msg, err)
	l.print(LvlErr, logMsg)

	l.sigChan <- syscall.SIGTERM
}

func (l *Logger) Error(msg string, err error) {
	logMsg := fmt.Sprintf("[ERROR]\t%s:\t%v", msg, err)
	l.print(LvlErr, logMsg)
}

func (l *Logger) Warning(msg string) {
	logMsg := fmt.Sprintf("[WARNING]\t%s", msg)
	l.print(LvlWrn, logMsg)
}

func (l *Logger) Info(msg string) {
	logMsg := fmt.Sprintf("[INFO]\t%s", msg)
	l.print(LvlInf, logMsg)
}

func (l *Logger) Debug(msg string) {
	logMsg := fmt.Sprintf("[DEBUG]\t%s", msg)
	l.print(LvlDeb, logMsg)
}
