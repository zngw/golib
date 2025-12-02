// Package zmap
// @Description 封装sync.Map, 添加一个源子计数，容易获取map数量
// @Author  55
// @Date  2025/12/2
package zmap

import (
	"sync"
	"sync/atomic"
)

type ZMap struct {
	m     sync.Map
	count atomic.Int32
}

// New 创建
func New() *ZMap {
	return &ZMap{
		m: sync.Map{},
	}
}

// Store 存储键值对，存在则替换
func (zm *ZMap) Store(key, value any) {
	_, _ = zm.Swap(key, value)
}

// Swap 替换m中的值，如果m中不存在，则新加。
func (zm *ZMap) Swap(key, value any) (previous any, loaded bool) {
	if previous, loaded = zm.m.Swap(key, value); !loaded {
		// 是新 key，计数 +1
		zm.count.Add(1)
	}

	return
}

// LoadOrStore 如果m中存在返回m中的值（不更新）；如果不存在，插入并返回插入值
func (zm *ZMap) LoadOrStore(key, value any) (actual any, loaded bool) {
	if actual, loaded = zm.m.LoadOrStore(key, value); !loaded {
		zm.count.Add(1)
	}

	return
}

// CompareAndSwap 如果m中存在key，且值等于old， 不更新值返回false;否则更新值，返回true
func (zm *ZMap) CompareAndSwap(key, old, new any) bool {
	return zm.m.CompareAndSwap(key, old, new)
}

// CompareAndDelete 如果m中存在key，且其值等于old，则删除该键值对并返回 true；否则返回 false
func (zm *ZMap) CompareAndDelete(key, old any) (deleted bool) {
	return zm.m.CompareAndDelete(key, old)
}

// LoadAndDelete 删除并返回删除的删
func (zm *ZMap) LoadAndDelete(key any) (value any, loaded bool) {
	if value, loaded = zm.m.LoadAndDelete(key); loaded {
		zm.count.Add(-1)
	}
	return
}

// Delete 删除键，并更新计数
func (zm *ZMap) Delete(key any) {
	zm.LoadAndDelete(key)
}

// Load 获取值
func (zm *ZMap) Load(key any) (value any, ok bool) {
	return zm.m.Load(key)
}

// Len 获取元素数量
func (zm *ZMap) Len() int {
	return int(zm.count.Load())
}

// Range 遍历
func (zm *ZMap) Range(f func(key, value any) bool) {
	zm.m.Range(f)
}
