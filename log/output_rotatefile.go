package log

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"
)

var defaultLogFileName = "file.log"

type rotateFileMode string

const (
	rotateFileModeNone  rotateFileMode = ""
	rotateFileModeDaily rotateFileMode = "daily"
)

var _ io.WriteCloser = (*rotateFileWriter)(nil)

type rotateFileConfig struct {
	FileName string
	Mode     rotateFileMode
	MaxDays  int
}

type rotateFileWriter struct {
	cfg rotateFileConfig

	mu   sync.Mutex
	file *os.File
	done chan struct{}
}

func newRotateFileWriter(cfg rotateFileConfig) *rotateFileWriter {
	if cfg.FileName == "" {
		cfg.FileName = defaultLogFileName
	}
	fw := &rotateFileWriter{
		cfg:  cfg,
		done: make(chan struct{}),
	}
	return fw
}

func (fw *rotateFileWriter) Init() {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	if fw.done != nil {
		close(fw.done)
	}
	fw.done = make(chan struct{})
	if fw.cfg.Mode == rotateFileModeDaily {
		go fw.dailyRotate()
	}
}

func (fw *rotateFileWriter) Write(p []byte) (n int, err error) {
	return fw.WriteLog(p, LevelInfo)
}

func (fw *rotateFileWriter) WriteLog(p []byte, _ Level) (int, error) {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	if fw.file == nil {
		if err := fw.openExistingOrNew(); err != nil {
			return 0, err
		}
	}

	n, err := fw.file.Write(p)
	return n, err
}

func (fw *rotateFileWriter) CloseLog() {
	err := fw.Close()
	if err != nil {
		return
	}
}

func (fw *rotateFileWriter) Rotate() error {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	return fw.rotate()
}

func (fw *rotateFileWriter) rotate() error {
	if err := fw.closeFile(); err != nil {
		return err
	}
	if err := fw.openNew(); err != nil {
		return err
	}
	_ = fw.clearFiles()
	return nil
}

func (fw *rotateFileWriter) openExistingOrNew() error {
	_, err := os.Stat(fw.cfg.FileName)
	if os.IsNotExist(err) {
		return fw.openNew()
	}
	if err != nil {
		return fmt.Errorf("get stat of logfile error: %s", err)
	}

	file, err := os.OpenFile(fw.cfg.FileName, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fw.openNew()
	}
	fw.file = file
	return nil
}

func (fw *rotateFileWriter) openNew() error {
	err := os.MkdirAll(fw.dir(), 0o755)
	if err != nil {
		return fmt.Errorf("mkdir directories [%s] for new logfile error: %s", fw.dir(), err)
	}

	mode := os.FileMode(0o600)
	info, err := os.Stat(fw.cfg.FileName)
	if err == nil {
		mode = info.Mode()
		newName := fw.backupName(fw.cfg.FileName, time.Now())
		if err := os.Rename(fw.cfg.FileName, newName); err != nil {
			return fmt.Errorf("rename logfile error: %s", err)
		}
	}

	f, err := os.OpenFile(fw.cfg.FileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("open new logfile error: %s", err)
	}
	fw.file = f
	return nil
}

var backupTimeFormat = "20060102-150405"

func (fw *rotateFileWriter) backupName(name string, t time.Time) string {
	dir := filepath.Dir(name)
	filename := filepath.Base(name)
	ext := filepath.Ext(filename)
	prefix := filename[:len(filename)-len(ext)]

	timestamp := t.Format(backupTimeFormat)
	return filepath.Join(dir, fmt.Sprintf("%s.%s%s", prefix, timestamp, ext))
}

func (fw *rotateFileWriter) parseTimeFromBackupName(filename, prefix, ext string) (time.Time, error) {
	if !strings.HasPrefix(filename, prefix) {
		return time.Time{}, errors.New("missmatched prefix")
	}
	if !strings.HasSuffix(filename, ext) {
		return time.Time{}, errors.New("missmatched ext")
	}
	if len(prefix) >= len(filename)-len(ext) {
		return time.Time{}, errors.New("missmatched prefix and ext")
	}
	timestamp := filename[len(prefix) : len(filename)-len(ext)]
	return time.ParseInLocation(backupTimeFormat, timestamp, time.Local)
}

func (fw *rotateFileWriter) dir() string {
	return filepath.Dir(fw.cfg.FileName)
}

func (fw *rotateFileWriter) dailyRotate() {
	fw.mu.Lock()
	doneCh := fw.done
	fw.mu.Unlock()
	for {
		now := time.Now()
		// Calculate the time difference until the next hour.
		nextHour := now.Truncate(time.Hour).Add(time.Hour)
		select {
		case <-time.After(nextHour.Sub(now)):
		case <-doneCh:
			return
		}

		// Rotate the log file at 0 hour of the day.
		if nextHour.Hour() == 0 {
			_ = fw.Rotate()
			// Ensure it's executed only once, even if the waiting period crosses midnight.
			time.Sleep(time.Minute)
		}
	}
}

type logFileInfo struct {
	info os.FileInfo
	t    time.Time
}

func (fw *rotateFileWriter) oldLogFiles() ([]logFileInfo, error) {
	entries, err := os.ReadDir(fw.dir())
	if err != nil {
		return nil, fmt.Errorf("read log file directory error: %s", err)
	}
	fileInfos := make([]logFileInfo, 0)

	filename := filepath.Base(fw.cfg.FileName)
	ext := filepath.Ext(filename)
	prefix := filename[:len(filename)-len(ext)] + "."

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}

		if t, err := fw.parseTimeFromBackupName(entry.Name(), prefix, ext); err == nil {
			fileInfos = append(fileInfos, logFileInfo{info: info, t: t})
			continue
		}
	}

	slices.SortFunc(fileInfos, func(a, b logFileInfo) int {
		return a.t.Compare(b.t)
	})
	return fileInfos, nil
}

func (fw *rotateFileWriter) clearFiles() error {
	if fw.cfg.Mode == rotateFileModeNone {
		return nil
	}
	if fw.cfg.Mode == rotateFileModeDaily && fw.cfg.MaxDays <= 0 {
		return nil
	}

	files, err := fw.oldLogFiles()
	if err != nil {
		return err
	}

	var toRemove []logFileInfo
	if fw.cfg.Mode == rotateFileModeDaily {
		cutoff := time.Now().Add(-time.Duration(fw.cfg.MaxDays) * time.Duration(24) * time.Hour).Add(5 * time.Millisecond)
		for _, f := range files {
			if f.t.Before(cutoff) {
				toRemove = append(toRemove, f)
			}
		}
	}

	for _, f := range toRemove {
		_ = os.Remove(filepath.Join(fw.dir(), f.info.Name()))
	}
	return nil
}

func (fw *rotateFileWriter) Close() error {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	if fw.done != nil {
		close(fw.done)
		fw.done = nil
	}
	return fw.closeFile()
}

func (fw *rotateFileWriter) closeFile() error {
	if fw.file == nil {
		return nil
	}
	err := fw.file.Close()
	fw.file = nil
	return err
}
