package log

import "fmt"

type Logger interface {
	Debug(msg string, ctx ...interface{})
	Info(msg string, ctx ...interface{})
	Warn(msg string, ctx ...interface{})
	Error(msg string, ctx ...interface{})
}

type logger struct{}

func New(ctx ...interface{}) Logger {
	return &logger{}
}

func (l *logger) Debug(msg string, ctx ...interface{}) {}

func (l *logger) Info(msg string, ctx ...interface{}) {}

func (l *logger) Warn(msg string, ctx ...interface{}) {}

func (l *logger) Error(msg string, ctx ...interface{}) {
	fmt.Println("Error:", msg, ctx)
}
