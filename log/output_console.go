package log

import (
	"io"
	"os"
	"strings"
)

// brush is a color join function
type brush func(string) string

// newBrush returns a fix color Brush
func newBrush(color string) brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

func emptyBrush(text string) string {
	return text
}

var colors = []brush{
	emptyBrush,       // Trace              No Color
	newBrush("1;36"), // Debug              Light Cyan
	newBrush("1;34"), // Info 				Blue
	newBrush("1;33"), // Warn               Yellow
	newBrush("1;31"), // Error              Red
}

func colorBrushByLevel(level Level) brush {
	switch level {
	case LevelTrace:
		return colors[0]
	case LevelDebug:
		return colors[1]
	case LevelInfo:
		return colors[2]
	case LevelWarn:
		return colors[3]
	case LevelError:
		return colors[4]
	default:
		return colors[2]
	}
}

var _ io.Writer = (*consoleWriter)(nil)

type consoleConfig struct {
	Colorful bool
}

type consoleWriter struct {
	cfg consoleConfig
	w   io.Writer
}

func newConsoleWriter(cfg consoleConfig) io.Writer {
	return &consoleWriter{
		cfg: cfg,
		w:   os.Stdout,
	}
}

func (cw *consoleWriter) Write(p []byte) (n int, err error) {
	return cw.w.Write(p)
}

func (cw *consoleWriter) WriteLog(p []byte, level Level) (n int, err error) {
	if cw.cfg.Colorful {
		p = []byte(strings.Replace(string(p), level.LogPrefix(), colorBrushByLevel(level)(level.LogPrefix()), 1))
	}

	return cw.w.Write(p)
}

func (cw *consoleWriter) CloseLog() {

}
