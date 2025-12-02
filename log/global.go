package log

var logger = New(Option{
	DisableLogColor: false,
	CallerSkip:      1,
})

func Init(opt Option) {
	logger.WithOptions(opt)
}

func Error(format string, v ...any) {
	logger.Error(format, v...)
}

func Warn(format string, v ...any) {
	logger.Warn(format, v...)
}

func Info(format string, v ...any) {
	logger.Info(format, v...)
}

func Debug(format string, v ...any) {
	logger.Debug(format, v...)
}

func Trace(format string, v ...any) {
	logger.Trace(format, v...)
}

func Log(level Level, offset int, msg string, args ...any) {
	logger.Log(level, offset, msg, args...)
}
