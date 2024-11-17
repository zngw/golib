package log

var logger = New(Option{
	DisableLogColor: false,
})

func Init(opt Option) {
	logger = logger.WithOptions(opt)
}

func Error(format string, v ...interface{}) {
	logger.log(LevelError, 0, format, v...)
}

func Warn(format string, v ...interface{}) {
	logger.log(LevelWarn, 0, format, v...)
}

func Info(format string, v ...interface{}) {
	logger.log(LevelInfo, 0, format, v...)
}

func Debug(format string, v ...interface{}) {
	logger.log(LevelDebug, 0, format, v...)
}

func Trace(format string, v ...interface{}) {
	logger.log(LevelTrace, 0, format, v...)
}

func Log(level Level, offset int, msg string, args ...interface{}) {
	logger.log(level, offset, msg, args...)
}
