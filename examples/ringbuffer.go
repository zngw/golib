// @Title
// @Description $
// @Author  55
// @Date  2022/5/30
package main

import (
	"fmt"
	"github.com/zngw/golib/ringbuffer"
)

func main() {
	// 初始化单元大小为4的RingBuffer
	rb, err := ringbuffer.NewRingBuffer(4)
	if err != nil {
		panic(err.Error())
	}

	// 写入9个元素
	for i := 0; i < 9; i++ {
		rb.Write(i)
		fmt.Printf("write data:%v, capacity:%v, len:%v \n", i, rb.Capacity(), rb.Len())
	}

	// 弹出5个元素
	for i := 0; i < 5; i++ {
		data := rb.Pop()
		fmt.Printf("pop data:%v, capacity:%v, len:%v \n", data, rb.Capacity(), rb.Len())
	}

	// 再写入8个
	for i := 0; i < 8; i++ {
		rb.Write(i)
		fmt.Printf("write data:%v, capacity:%v, len:%v \n", i, rb.Capacity(), rb.Len())
	}

	// 弹出所有元素
	for !rb.IsEmpty() {
		data := rb.Pop()
		fmt.Printf("pop data:%v, capacity:%v, len:%v \n", data, rb.Capacity(), rb.Len())
	}
}
