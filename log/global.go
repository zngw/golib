package log

var logger = New(Option{
	DisableLogColor: false,
	CallerSkip:      1,
})

func Init(opt Option) {
	logger.WithOptions(opt)
}

func Error(format string, v ...interface{}) {
	logger.Error(format, v...)
}

func Warn(format string, v ...interface{}) {
	logger.Warn(format, v...)
}

func Info(format string, v ...interface{}) {
	logger.Info(format, v...)
}

func Debug(format string, v ...interface{}) {
	logger.Debug(format, v...)
}

func Trace(format string, v ...interface{}) {
	logger.Trace(format, v...)
}

func Log(level Level, offset int, msg string, args ...interface{}) {
	logger.Log(level, offset, msg, args...)
}
