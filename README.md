# golib

一些通用的go lib 封装

| 模块 | 包 | 描述                                   |
|-----|----|--------------------------------------|
| 日志 | log | 提供终端输出、文件输出的日志，日志分level等级、tags来控制输出  |
| 双向环形链表 | ringbuffer | |
| 无限缓存chan | zchan | 一个使用双向环形链表，能无限缓存的chan通道 |
| 并发安全Set集合 | set | 一个基于sync.Map封装的线程安全Set集合 |
| 并发安全Map | zmap | 封装sync.Map, 添加一个源子计数，容易获取map数量|
| 并发安全切片| zslice | 封装一个可线程安全的Slice，未使用分段锁， 10k量级，超出性能急降|
|字符串与数字转换| str | 一些字符串与数字转换方法 |
| 加密算法封装 | crypt | 集成md5、hmac Hash算法，aes、des对称加密，rsa非对称加密 |

# 安装

```bash
go get -u github.com/zngw/golib
```

# 使用

### 日志

* 无初始化直接使用时，默认为终端打印日志

使用案例参考 [https://github.com/zngw/golib/blob/main/examples/log.go](https://github.com/zngw/golib/blob/main/examples/log.go)

### 双向环形链表

* 双向环形链表buf其结构类似于一个手串，手串上的珠子就可以当做是一个节点，每个节点可以是一个固定大小的数组
* 双向环形链表buf上分别有两个读写指针readCell和writeCell，指向将要进行读写操作的cell，负责进行数据读写
* readCell永远追赶writeCell，当追上时，代表写满了，进行扩容操作
* 扩容操作即在写指针的后面插入一个新建的空闲cell
* 缩容操作修改链表指向即可，让buf恢复原样，仅保持两个cell即可，其他cell由于不再被引用，会被GC自动回收
* 在链表写入(Write)和读取(Read)时用原子操作修改链表有效数据长度count

使用案例参考 [https://github.com/zngw/golib/blob/main/examples/ringbuffer.go](https://github.com/zngw/golib/blob/main/examples/ringbuffer.go)

### 无限缓存channel

* 一个使用双向环形链表，能无限缓存的chan通道

使用案例参考 [https://github.com/zngw/golib/blob/main/examples/zchan.go](https://github.com/zngw/golib/blob/main/examples/zchan.go)

### 线程安全Set集合

* 一个基于sync.Map封装的线程安全Set集合

使用案例参考 [https://github.com/zngw/golib/blob/main/examples/set.go](https://github.com/zngw/golib/blob/main/examples/set.go)

### 并发安全Map

* 封装sync.Map, 添加一个源子计数，容易获取map数量，其他接口如 sync.Map一样

### 并发安全切片

* 封装一个可线程安全的Slice
* 提供常见操作方法，如 Append, Get, Set, Len, Delete, Range, Clear, Sort, Find等
* 未使用分段锁，如果 slice 很大会影响效率
* 部分实测数据参考（8核8GB云服环境测试）:
* 1) 小于1k元素，100并发写，几乎无影响；
* 2) 10k元素左右时，100并发写，QPS ≈ 5k~10k
* 3) 100k元素，100并发写， QPS 骤降至 500~1k（因扩容+锁竞争），不建议使用

使用案例参考 [https://github.com/zngw/golib/blob/main/examples/zslice.go](https://github.com/zngw/golib/blob/main/examples/zslice.go)

### 加密算法封装

* Hash算法：md5、hmac。md5已不安全，建议使用hmac替代
* 对称加密：aes、des，des也不够安全了，建议使用aes
[https://github.com/zngw/golib/blob/main/examples/aes.go](https://github.com/zngw/golib/blob/main/examples/aes.go)
* 非对称加密：rsa， rsa不适合加密长字符串
[https://github.com/zngw/golib/blob/main/examples/rsa.go](https://github.com/zngw/golib/blob/main/examples/rsa.go)
