package logger

import (
	"log"
	"os"
)

type Logger interface {
	Info(a ...interface{})
	Debug(a ...interface{})

	Infof(format string, a ...interface{})
}

type logger struct {
	log     *log.Logger
	isDebug bool
}

func InitLog() *logger {
	return &logger{
		log:     log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime),
		isDebug: true,
	}
}
func (l *logger) Infof(format string, a ...interface{}) {
	l.log.Printf(format, a...)
}

func (l *logger) Info(a ...interface{}) {
	l.log.Println(a...)
}

func (l *logger) Debug(a ...interface{}) {
	if l.isDebug {
		l.log.Println(a...)
	}
}
