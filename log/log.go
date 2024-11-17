package log

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

type writer interface {
	WriteLog([]byte, Level) (n int, err error)
	CloseLog()
}

var defaultWriter = os.Stdout

type Logger struct {
	outMu sync.Mutex
	out   io.Writer

	level         Level
	tags          tags
	colorful      bool
	callerEnabled bool
	callerSkip    int
}

// New 新建日志对象， 使用`opt ...`的目的是为了让New可以缺省参数使用，实际只使用到了opt[0]
func New(opt ...Option) *Logger {
	l := &Logger{}

	if len(opt) > 0 {
		l.WithOptions(opt[0])
	}

	if l.out == nil {
		l.out = defaultWriter
	}
	if l.level == 0 {
		l.level = LevelTrace
	}
	return l
}

// WithOptions 修改当前的日志配置
func (l *Logger) WithOptions(opt Option) {
	if l.out == nil || opt.LogPath != "" {
		if l.out != nil {
			if lw, ok := l.out.(writer); ok {
				lw.CloseLog()
			}
		}

		if opt.LogPath == "" || opt.LogPath == "console" {
			l.out = newConsoleWriter(consoleConfig{
				Colorful: !opt.DisableLogColor,
			})
		} else {
			writer := newRotateFileWriter(rotateFileConfig{
				FileName: opt.LogPath,
				Mode:     rotateFileModeDaily,
				MaxDays:  opt.MaxDays,
			})
			writer.Init()
			l.out = writer
		}
	}

	if opt.LogLevel != "" {
		level, err := parseLevel(opt.LogLevel)
		if err == nil {
			l.level = level
		}
	}

	l.tags = parseTags(opt.Tags)
	l.callerEnabled = !opt.DisableCaller
	if opt.CallerSkip > 0 {
		l.callerSkip = opt.CallerSkip
	}
	return
}

func (l *Logger) clone() *Logger {
	clone := &Logger{
		out:           l.out,
		level:         l.level,
		callerEnabled: l.callerEnabled,
		callerSkip:    l.callerSkip,
	}
	return clone
}

func (l *Logger) log(level Level, offset int, msg string, args ...interface{}) {
	show := l.level.Enabled(level)
	if !show {
		return
	}

	tag := ""
	if len(args) > 0 && strings.Count(msg, "%")-strings.Count(msg, "%%")*2 == 0 {
		// 大于一个参数，有占位符
		tag = msg
		msg = args[0].(string)
		args = args[1:]
	}

	tag, show = l.tags.GetTag(tag)
	if !show {
		return
	} else {
		msg = tag + msg
	}

	caller := ""
	if l.callerEnabled {
		caller = getCallerPrefix(3 + l.callerSkip + offset)
	}

	outMsg := fmt.Sprintf("%s%s%s%s\n",
		time.Now().Format("2006-01-02 15:04:05.000 "),
		level.LogPrefix(),
		caller,
		getMessage(msg, args),
	)

	if lw, ok := l.out.(writer); ok {
		l.outMu.Lock()
		defer l.outMu.Unlock()
		_, _ = lw.WriteLog([]byte(outMsg), level)
		return
	}

	l.outMu.Lock()
	defer l.outMu.Unlock()
	_, _ = l.out.Write([]byte(outMsg))
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.log(LevelError, 0, format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.log(LevelWarn, 0, format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.log(LevelInfo, 0, format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.log(LevelDebug, 0, format, v...)
}

func (l *Logger) Trace(format string, v ...interface{}) {
	l.log(LevelTrace, 0, format, v...)
}

func (l *Logger) Log(level Level, offset int, msg string, args ...interface{}) {
	l.log(level, offset, msg, args...)
}

func getMessage(template string, fmtArgs []interface{}) string {
	if len(fmtArgs) == 0 {
		return template
	}

	if template != "" {
		return fmt.Sprintf(template, fmtArgs...)
	}

	if len(fmtArgs) == 1 {
		if str, ok := fmtArgs[0].(string); ok {
			return str
		}
	}
	return fmt.Sprint(fmtArgs...)
}

func getCallerPrefix(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "???"
		line = 0
	}
	_, file = path.Split(file)
	return fmt.Sprintf("[%s:%d] ", file, line)
}
