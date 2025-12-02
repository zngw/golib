// Package set
// @Description 基于sync.Map封装的线程安全Set集合
// @Author  55
// @Date  2022/5/30
package set

import (
	"sync"
	"sync/atomic"
)

// Set 集合
type Set struct {
	m   sync.Map
	num atomic.Int32
}

// New 创建
func New() *Set {
	return &Set{
		m: sync.Map{},
	}
}

// Add 添加
func (s *Set) Add(item any) {
	_, loaded := s.m.Swap(item, true)
	if !loaded {
		s.num.Add(1)
	}
}

// Remove 删除
func (s *Set) Remove(item any) {
	_, loaded := s.m.LoadAndDelete(item)
	if loaded {
		s.num.Add(-1)
	}
}

// Has 判断是否存在
func (s *Set) Has(item any) (ok bool) {
	_, ok = s.m.Load(item)
	return
}

// Len 获取集合大小
func (s *Set) Len() int {
	return int(s.num.Load())
}

// Clear 清除
func (s *Set) Clear() {
	s.m.Range(func(k, v any) bool {
		s.Remove(k)
		return true
	})
}

// IsEmpty 判断是否为空
func (s *Set) IsEmpty() bool {
	return s.num.Load() == 0
}

// List 转切片输出
func (s *Set) List() (list []any) {
	s.m.Range(func(k, v any) bool {
		list = append(list, k)
		return true
	})

	return list
}

// Range 遍历
func (s *Set) Range(f func(key any) bool) {
	s.m.Range(func(k, v any) bool {
		return f(k)
	})
}
