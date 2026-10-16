// Package zchan
// @Description 无限缓存channel
// @Author  55
// @Date  2022/5/30
package zchan

import (
	"fmt"

	"github.com/zngw/golib/ringbuffer"
)

type ZChan[T any] struct {
	In     chan<- T                  // 写入channel
	Out    <-chan T                  // 读取channel
	buffer *ringbuffer.RingBuffer[T] // 双向环形链表
}

// Len uc中总共的元素数量
func (uc *ZChan[T]) Len() int {
	return len(uc.In) + uc.BufLen() + len(uc.Out)
}

// BufLen uc的buf中的元素数量
func (uc *ZChan[T]) BufLen() int {
	return uc.buffer.Len()
}

// New 新建一个无限缓存的Channel，并指定In和Out大小(In和Out设置得一样大)
func New[T any](initCapacity int) (ch *ZChan[T], err error) {
	if initCapacity <= 0 {
		return nil, fmt.Errorf("initCapacity 必须大于 0")
	}

	rb, err := ringbuffer.NewRingBuffer[T](512)
	if err != nil {
		return
	}

	in := make(chan T, initCapacity)
	out := make(chan T, initCapacity)
	ch = &ZChan[T]{In: in, Out: out, buffer: rb}

	go process(in, out, ch)

	return
}

// 内部Worker Goroutine实现
func process[T any](in, out chan T, ch *ZChan[T]) {
	defer func() {
		drainBuffer(ch, out)
		close(out)
	}()

	// 不断从in中读取数据放入到out或者buf中
loop:
	for {
		// 第一步：从in中读取数据
		value, ok := <-in
		if !ok {
			// in 关闭了，退出loop
			break loop
		}

		// 第二步：将数据存储到out或者buf中
		if ch.buffer.Len() > 0 {
			// 当buf中有数据时，新数据优先存放到buf中，确保数据FIFO原则
			ch.buffer.Write(value)
		} else {
			// out 没有满,数据放入out中
			select {
			case out <- value:
				continue
			default:
			}

			// out 满了，数据放入buf中
			ch.buffer.Write(value)
		}

		// 第三步：处理buf，一直尝试把buf中的数据放入到out中，直到buf中没有数据
		for ch.buffer.Len() > 0 {
			// 先尝试从 in 收新数据，避免阻塞
			select {
			case val, ok := <-in:
				if !ok {
					break loop
				}
				ch.buffer.Write(val)
				continue // 继续尝试排 buffer
			default:
			}

			// in 当前没数据，尝试把 buffer 头部发到 out
			val, ok := ch.buffer.Peek()
			if !ok {
				break
			}
			select {
			case out <- val:
				ch.buffer.Pop() // 发送成功后再弹
				if ch.buffer.IsEmpty() {
					ch.buffer.Reset()
				}
			default:
				// out 也满了，暂时排不动，回到外层继续从 in 收
				break
			}
		}
	}
}

func drainBuffer[T any](ch *ZChan[T], out chan T) {
	for ch.buffer.Len() > 0 {
		val, ok := ch.buffer.Pop()
		if !ok {
			break
		}

		select {
		case out <- val:
		default:
			// out 已满或无人消费时直接丢弃，避免永久阻塞
		}
	}
}
