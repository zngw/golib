package log

// An Option configures a Logger.
type Option struct {
	LogPath         string // 日志输入文件， console为终端输出
	LogLevel        string // 日志等级
	Tags            string // 日志Tag
	MaxDays         int    // 日志文件保留日期
	DisableLogColor bool   // 终端输出是否显示颜色
	DisableCaller   bool   // 是否打印调用文件
	CallerSkip      int    // 打印文件级
}

func (opt *Option) apply(log *Logger) {
	if log.out == nil || opt.LogPath != "" {
		if log.out != nil {
			if lw, ok := log.out.(writer); ok {
				lw.CloseLog()
			}
		}

		if opt.LogPath == "" || opt.LogPath == "console" {
			log.out = newConsoleWriter(consoleConfig{
				Colorful: !opt.DisableLogColor,
			})
		} else {
			writer := newRotateFileWriter(rotateFileConfig{
				FileName: opt.LogPath,
				Mode:     rotateFileModeDaily,
				MaxDays:  opt.MaxDays,
			})
			writer.Init()
			log.out = writer
		}
	}

	if opt.LogLevel != "" {
		level, err := parseLevel(opt.LogLevel)
		if err == nil {
			log.level = level
		}
	}

	log.tags = parseTags(opt.Tags)
	log.callerEnabled = !opt.DisableCaller
	if opt.CallerSkip > 0 {
		log.callerSkip = opt.CallerSkip
	}
}
