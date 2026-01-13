// Package zmap
// @Description 封装sync.Map, 添加一个源子计数，容易获取map数量
// @Author  55
// @Date  2025/12/2
package zmap

import (
	"reflect"
	"sync"
)

// ZMap 是一个线程安全的map，针对高频读取和低频写入进行了优化。
// map keys < 1000时，其性能优于 sync.Map。
type ZMap[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

// New 创建一个新的 ZMap。
func New[K comparable, V any]() *ZMap[K, V] {
	return &ZMap[K, V]{
		data: make(map[K]V),
	}
}

// Store 存储键值对
func (zm *ZMap[K, V]) Store(key K, value V) {
	zm.mu.Lock()
	zm.data[key] = value
	zm.mu.Unlock()
}

// LoadOrStore 如果键值已存在，LoadOrStore 函数返回该键的现有值。
// 否则，它会存储并返回给定的值。
// 如果值已加载，则 loaded 结果为 true；如果值已存储，则 loaded 结果为 false。
func (zm *ZMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	zm.mu.RLock()
	if v, ok := zm.data[key]; ok {
		zm.mu.RUnlock()
		return v, true
	}
	zm.mu.RUnlock()

	// 仔细检查写锁内部（以避免并发存储期间的竞争条件）
	zm.mu.Lock()
	defer zm.mu.Unlock()
	if v, ok := zm.data[key]; ok {
		return v, true
	}
	zm.data[key] = value
	return value, false
}

// CompareAndSwap 如果m中存在key，且值等于old， 不更新值返回false;否则更新值，返回true
func (zm *ZMap[K, V]) CompareAndSwap(key K, old, new V) bool {
	zm.mu.Lock()
	defer zm.mu.Unlock()

	current, exists := zm.data[key]
	if !exists {
		return false
	}

	// 使用泛型"深度值相等"来比较，不是指针相等
	// 性能比 == 慢约 10～100 倍（取决于数据结构复杂度）
	if !reflect.DeepEqual(current, old) {
		return false
	}

	zm.data[key] = new
	return true
}

// CompareAndDelete 如果m中存在key，且其值等于old，则删除该键值对并返回 true；否则返回 false
func (zm *ZMap[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	zm.mu.Lock()
	defer zm.mu.Unlock()

	current, exists := zm.data[key]
	if !exists {
		return false
	}

	// 使用泛型"深度值相等"来比较，不是指针相等
	// 性能比 == 慢约 10～100 倍（取决于数据结构复杂度）
	if !reflect.DeepEqual(current, old) {
		return false
	}

	delete(zm.data, key)
	return true
}

// LoadAndDelete 删除并返回删除的删
func (zm *ZMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	zm.mu.Lock()
	defer zm.mu.Unlock()

	value, loaded = zm.data[key]
	if loaded {
		delete(zm.data, key)
	}
	return
}

// Delete 删除键
func (zm *ZMap[K, V]) Delete(key K) {
	zm.mu.Lock()
	delete(zm.data, key)
	zm.mu.Unlock()
}

// Load 函数返回map中存储的键对应的值，如果不存在则返回零值和 false。
func (zm *ZMap[K, V]) Load(key K) (value V, ok bool) {
	zm.mu.RLock()
	value, ok = zm.data[key]
	zm.mu.RUnlock()
	return
}

// Len 获取元素数量
func (zm *ZMap[K, V]) Len() int {
	zm.mu.RLock()
	defer zm.mu.RUnlock()
	return len(zm.data)
}

// Range 会依次对map中的每个键和值调用 f 函数。
// 如果 f 返回 false，range 将停止迭代。
// 注意：在整个迭代过程中，Range 会持有读取锁！
// 对于大型映射表或 f 函数执行速度较慢的情况，请考虑先复制键/值对。
func (zm *ZMap[K, V]) Range(f func(key K, value V) bool) {
	zm.mu.RLock()
	defer zm.mu.RUnlock()
	for k, v := range zm.data {
		if !f(k, v) {
			break
		}
	}
}

// Keys 返回map中所有键的快照。
func (zm *ZMap[K, V]) Keys() []K {
	zm.mu.RLock()
	defer zm.mu.RUnlock()
	keys := make([]K, 0, len(zm.data))
	for k := range zm.data {
		keys = append(keys, k)
	}
	return keys
}

// Values 返回map中所有值的快照。
func (zm *ZMap[K, V]) Values() []V {
	zm.mu.RLock()
	defer zm.mu.RUnlock()
	values := make([]V, 0, len(zm.data))
	for _, v := range zm.data {
		values = append(values, v)
	}
	return values
}

// Clone 返回整个map的浅拷贝
func (zm *ZMap[K, V]) Clone() map[K]V {
	zm.mu.RLock()
	defer zm.mu.RUnlock()
	clone := make(map[K]V, len(zm.data))
	for k, v := range zm.data {
		clone[k] = v
	}
	return clone
}
