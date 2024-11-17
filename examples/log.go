package main

import "github.com/zngw/golib/log"

func main() {
	// 无初始化默认终端输出

	// 输出到文件日志配置
	log.Init(log.Option{
		LogPath:         "log/file.log", // 日志文件路径。缺省为终端输出
		LogLevel:        "trace",        // 输出等级，缺省为trace，可选 trace、info、debug、warn、error
		Tags:            "",             // Tag, 缺省为显示所有tag，调用输出不带tag时，不受tag标签影响
		MaxDays:         7,              // 日志文件保留天数，仅在文件模式下生效，缺省为永久保留
		DisableLogColor: false,          // 是否禁用日志颜色显示，仅在终端模式下生效，缺省为不禁用
		DisableCaller:   false,          // 是否禁用显示打印所在文件及行数，缺省为不禁用
		CallerSkip:      0,              // 打印日志文件调用层级参数，缺省为0，即当前掉用log.Trace接口所在文件行数
	})

	// format 传入一个字符串
	log.Trace("直接输出")

	// format 传入带'%'的格式化字符串，且后面存在对应参数
	log.Info("格式化输出 %10s", "符字串")
	log.Debug("格式化输出 %10s， %d", "符字串", 55)

	// format 传入不带'%'格式化的tag字符串，且后面存在一个参数，format为tag，v参数为输出日志
	log.Warn("sys", "带tag输出")

	// format 传入不带'%'格式化的tag字符串，v[0]为带'%'的格式化字符串,后面存在对应参数，format为tag，v为格式化字符串数组
	log.Error("sys", "带tag格式化输出 %10s", "符字串")
	log.Log(log.LevelInfo, 0, "net", "带tag格式化输出 %10s， %d", "符字串", 55)

	// 可以创建多个独立的日志对象单独输出
	errlog := log.New(log.Option{
		LogPath:  "log/err.log",
		LogLevel: "error",
		MaxDays:  7,
	})
	errlog.Error("sys", "这是一个错误日志输出")

	mylog := log.New(log.Option{
		LogPath: "console",
	})
	mylog.Trace("net", "mylog 日志输出")
}
