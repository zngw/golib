package log

// Option 日志配置
type Option struct {
	LogPath         string // 日志输入文件， console为终端输出
	LogLevel        string // 日志等级
	Tags            string // 日志Tag
	MaxDays         int    // 日志文件保留日期
	DisableLogColor bool   // 终端输出是否显示颜色
	DisableCaller   bool   // 是否打印调用文件
	CallerSkip      int    // 打印文件级
}
