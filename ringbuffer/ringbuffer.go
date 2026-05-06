// Package ringbuffer
// @Description $ 双向环形链表
// @Author  55
// @Date  2022/5/30
package ringbuffer

import (
	"errors"
	"fmt"
	"sync"
)

var ErrIsEmpty = errors.New("ringbuffer is empty")

type cell[T any] struct {
	Data     []T      // 数据部分
	fullFlag bool     // cell满的标志
	next     *cell[T] // 指向后一个cellBuffer
	pre      *cell[T] // 指向前一个cellBuffer

	r int // 下一个要读的指针
	w int // 下一个要下的指针
}

type RingBuffer[T any] struct {
	mu        sync.RWMutex
	cellSize  int // cell大小
	cellCount int // cell数量
	count     int // 有效元素个数

	readCell  *cell[T] // 下一个要读的cell
	writeCell *cell[T] // 下一个要写的cell
}

// NewRingBuffer 新建一个RingBuffer，包含两个cell
func NewRingBuffer[T any](cellSize int) (buf *RingBuffer[T], err error) {
	if cellSize <= 0 || cellSize&(cellSize-1) != 0 {
		err = fmt.Errorf("初始大小必须是 2 的幂")
		return
	}

	rootCell := &cell[T]{
		Data: make([]T, cellSize),
	}
	lastCell := &cell[T]{
		Data: make([]T, cellSize),
	}
	rootCell.pre = lastCell
	lastCell.pre = rootCell
	rootCell.next = lastCell
	lastCell.next = rootCell

	buf = &RingBuffer[T]{
		cellSize:  cellSize,
		cellCount: 2,
		// count:   不需要显性初始化，默认为0
		readCell:  rootCell,
		writeCell: rootCell,
	}

	return
}

// Read 读取数据
func (r *RingBuffer[T]) Read() (data T, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.count == 0 {
		err = ErrIsEmpty
		return
	}

	// 读取数据，并将读指针向右移动一位
	data = r.readCell.Data[r.readCell.r]
	r.readCell.r++
	r.count--

	// 此cell已经读完
	if r.readCell.r == r.cellSize {
		// 读指针归零，并将该cell状态置为非满
		r.readCell.r = 0
		r.readCell.fullFlag = false
		// 将readCell指向下一个cell
		r.readCell = r.readCell.next
	}

	return
}

// Pop 读一个元素，读完后移动指针
func (r *RingBuffer[T]) Pop() (data T) {
	data, err := r.Read()
	if errors.Is(err, ErrIsEmpty) {
		fmt.Print(ErrIsEmpty.Error())
		return
	}
	return
}

// Peek 窥视 读一个元素，仅读但不移动指针
func (r *RingBuffer[T]) Peek() (data T) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.count == 0 {
		fmt.Print(ErrIsEmpty.Error())
		return
	}

	// 仅读
	data = r.readCell.Data[r.readCell.r]
	return
}

// Write 写入数据
func (r *RingBuffer[T]) Write(value T) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 在 r.writeCell.w 位置写入数据，指针向右移动一位
	r.writeCell.Data[r.writeCell.w] = value
	r.writeCell.w++
	r.count++

	// 当前cell写满了
	if r.writeCell.w == r.cellSize {
		// 指针置0，将该cell标记为已满，并指向下一个cell
		r.writeCell.w = 0
		r.writeCell.fullFlag = true
		r.writeCell = r.writeCell.next
	}

	// 下一个cell也已满，扩容
	if r.writeCell.fullFlag {
		r.grow()
	}
}

// grow 扩容（调用方需持有锁）
func (r *RingBuffer[T]) grow() {
	// 新建一个cell
	newCell := &cell[T]{
		Data: make([]T, r.cellSize),
	}

	// 总共三个cell，writeCell，preCell，newCell
	// 本来关系： preCell <===> writeCell
	// 现在将newcell插入：preCell <===> newCell <===> writeCell
	pre := r.writeCell.pre
	pre.next = newCell
	newCell.pre = pre
	newCell.next = r.writeCell
	r.writeCell.pre = newCell

	// 将writeCell指向新建的cell
	r.writeCell = r.writeCell.pre

	// cell 数量加一
	r.cellCount++
}

// IsEmpty 判断RingBuffer是否为空
func (r *RingBuffer[T]) IsEmpty() bool {
	return r.Len() == 0
}

// Capacity RingBuffer容量
func (r *RingBuffer[T]) Capacity() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.cellCount * r.cellSize
}

// Len RingBuffer数据长度
func (r *RingBuffer[T]) Len() (count int) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count = r.count
	return
}

// Reset 重置为仅指向两个cell的ring
func (r *RingBuffer[T]) Reset() {
	// 没有数据切cellCount只有两个时，无需重置
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.count == 0 && r.cellCount == 2 {
		return
	}

	// 保存两个保留的cell
	keepFirst := r.readCell
	keepLast := r.readCell.next

	// 断开中间cell的引用，帮助GC回收
	// 从keepLast.next开始遍历，直到回到keepFirst
	cur := keepLast.next
	for cur != keepFirst {
		next := cur.next
		cur.next = nil
		cur.pre = nil
		cur = next
	}

	keepLast.w = 0
	keepLast.r = 0
	keepFirst.r = 0
	keepFirst.w = 0
	r.cellCount = 2
	r.count = 0

	keepLast.next = keepFirst
	keepFirst.pre = keepLast
}
