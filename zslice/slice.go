// Package zslice
// @Description 封装一个可线程安全的Slice
// 提供常见操作方法，如 Append, Get, Set, Len, Delete, Range 等
// 未使用分段锁，如果 slice 很大会影响效率
// 部分实测数据参考（8核8GB云服环境测试）:
//	1) 小于1k元素，100并发写，几乎无影响；
//  2) 10k元素左右时，100并发写，QPS ≈ 5k~10k
// 	3) 100k元素，100并发写， QPS 骤降至 500~1k（因扩容+锁竞争），不建议使用
// @Author  55
// @Date  2025/12/28
package zslice

import (
	"sort"
	"sync"
)

// Slice 是一个线程安全的 slice 封装
type Slice[T any] struct {
	data []T
	mu   sync.RWMutex
}

// NewSlice 创建一个新的线程安全 slice
func NewSlice[T any]() *Slice[T] {
	return &Slice[T]{
		data: make([]T, 0),
	}
}

// Append 向 slice 末尾追加元素
func (s *Slice[T]) Append(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = append(s.data, item)
}

// Get 获取指定索引的元素
func (s *Slice[T]) Get(index int) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if index < 0 || index >= len(s.data) {
		var zero T
		return zero, false
	}
	return s.data[index], true
}

// Set 设置指定索引的元素
func (s *Slice[T]) Set(index int, item T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if index < 0 || index >= len(s.data) {
		return false
	}
	s.data[index] = item
	return true
}

// Len 获取 slice 长度
func (s *Slice[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

// Delete 删除指定索引的元素
func (s *Slice[T]) Delete(index int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if index < 0 || index >= len(s.data) {
		return false
	}
	s.data = append(s.data[:index], s.data[index+1:]...)
	return true
}

// DeleteBy 删除指定值（自定义相等判断）
func (s *Slice[T]) DeleteBy(equal func(T) bool) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	count := 0
	for i := len(s.data) - 1; i >= 0; i-- {
		if equal(s.data[i]) {
			s.data = append(s.data[:i], s.data[i+1:]...)
			count++
		}
	}
	return count
}

// Range 遍历 slice（类似 map.Range）
func (s *Slice[T]) Range(f func(index int, item T) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i, item := range s.data {
		if !f(i, item) {
			break
		}
	}
}

// ToSlice 返回一个当前 slice 的副本（非线程安全）
func (s *Slice[T]) ToSlice() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]T, len(s.data))
	copy(result, s.data)
	return result
}

// Sort 按自定义比较函数排序
func (s *Slice[T]) Sort(less func(T, T) bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sort.Slice(s.data, func(i, j int) bool {
		return less(s.data[i], s.data[j])
	})
}

// Find 查找满足条件的第一个元素
func (s *Slice[T]) Find(predicate func(T) bool) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, item := range s.data {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// FindAll 查找所有满足条件的元素
func (s *Slice[T]) FindAll(predicate func(T) bool) []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]T, 0)
	for _, item := range s.data {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Count 统计满足条件的元素个数
func (s *Slice[T]) Count(predicate func(T) bool) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	count := 0
	for _, item := range s.data {
		if predicate(item) {
			count++
		}
	}
	return count
}
