package logger

import (
	"fmt"
	"log"
	"os"
)

// "github.com/ynuraddi/t-tsarka/ilogger"

// var _ ilogger.ILogger = (*Logger)(nil)

const (
	errLvl = iota
	wrnLvl
	infLvl
	debLvl
)

type Logger struct {
	file  *os.File
	level int
}

func NewLogger(dst string, logLevel int) (*Logger, error) {
	file, err := os.OpenFile(dst, os.O_WRONLY, 0o666)
	if err != nil {
		return nil, err
	}

	return &Logger{
		file: file,
	}, nil
}

func (l *Logger) Close() {
	l.file.Close()
}

func (l *Logger) writeLog(logStr string) {
	log.Println(logStr)
	if l.file != nil {
		l.file.WriteString(logStr + "\n")
	}
}

func (l *Logger) Error(msg string, err error) {
	logStr := fmt.Sprintf("[ERROR] %s: %v", msg, err)
	if l.level <= errLvl {
		l.writeLog(logStr)
	}
}

func (l *Logger) Warning(msg string, err error) {
	logStr := fmt.Sprintf("[WARNING] %s: %v", msg, err)
	if l.level <= wrnLvl {
		l.writeLog(logStr)
	}
}

func (l *Logger) Info(msg string, err error) {
	logStr := fmt.Sprintf("[INFO] %s: %v", msg, err)
	if l.level <= infLvl {
		l.writeLog(logStr)
	}
}

func (l *Logger) Debug(msg string, err error) {
	logStr := fmt.Sprintf("[DEBUG] %s: %v", msg, err)
	if l.level <= debLvl {
		l.writeLog(logStr)
	}
}
