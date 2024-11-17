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

func New(opts ...Option) *Logger {
	l := &Logger{}

	for _, opt := range opts {
		opt.apply(l)
	}

	if l.out == nil {
		l.out = defaultWriter
	}
	if l.level == 0 {
		l.level = LevelTrace
	}
	return l
}

// WithOptions returns a new Logger with the given Options.
// It does not modify the original Logger.
func (l *Logger) WithOptions(opts ...Option) *Logger {
	c := l.clone()
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
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
