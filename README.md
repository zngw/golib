# golib

一些通用的go lib

# log - 日志

## 安装

```bash
go get -u github.com/zngw/golib/log
```

## 初始化

* 无初始化直接使用时，默认为终端打印日志

```go
// 初始化，输出到文件日志
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
```

使用案例参考 [https://github.com/zngw/golib/blob/main/examples/log.go](https://github.com/zngw/golib/blob/main/examples/log.go)

# set - Set集合

一个基于sync.Map封装的线程安全Set集合

## 安装使用

```go
go get -u github.com/zngw/golib/set
```

使用案例参考 [https://github.com/zngw/golib/blob/main/examples/set.go](https://github.com/zngw/golib/blob/main/examples/set.go)

# ringbuffer - 双向环形链表

* 双向环形链表buf其结构类似于一个手串，手串上的珠子就可以当做是一个节点，每个节点可以是一个固定大小的数组
* 双向环形链表buf上分别有两个读写指针readCell和writeCell，指向将要进行读写操作的cell，负责进行数据读写
* readCell永远追赶writeCell，当追上时，代表写满了，进行扩容操作
* 扩容操作即在写指针的后面插入一个新建的空闲cell
* 缩容操作修改链表指向即可，让buf恢复原样，仅保持两个cell即可，其他cell由于不再被引用，会被GC自动回收
* 在链表写入(Write)和读取(Read)时用原子操作修改链表有效数据长度count

## 安装使用

```go
go get -u github.com/zngw/golib/ringbuffer
```

使用案例参考 [https://github.com/zngw/golib/blob/main/examples/ringbuffer.go](https://github.com/zngw/golib/blob/main/examples/ringbuffer.go)

# zchan - 无限缓存channel

安装使用

```go
go get -u github.com/zngw/golib/zchan
```

使用案例参考 [https://github.com/zngw/golib/blob/main/examples/zchan.go](https://github.com/zngw/golib/blob/main/examples/zchan.go)

