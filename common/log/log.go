package log

var root = logger{}

func Debug(msg string, ctx ...interface{}) {
	root.Debug(msg, ctx)
}

func Info(msg string, ctx ...interface{}) {
	root.Info(msg, ctx)
}

func Warn(msg string, ctx ...interface{}) {
	root.Warn(msg, ctx)
}

func Error(msg string, ctx ...interface{}) {
	root.Error(msg, ctx)
}
