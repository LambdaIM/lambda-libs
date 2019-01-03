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

func (l *logger) Debug(msg string, ctx ...interface{}) {
	fmt.Println("Debug:", msg, ctx)
}

func (l *logger) Info(msg string, ctx ...interface{}) {
	fmt.Println("Info:", msg, ctx)
}

func (l *logger) Warn(msg string, ctx ...interface{}) {
	fmt.Println("Warn:", msg, ctx)
}

func (l *logger) Error(msg string, ctx ...interface{}) {
	fmt.Println("Error:", msg, ctx)
}
