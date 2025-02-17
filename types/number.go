package types

import (
	"sync/atomic"

	"gopkg.in/mgo.v2/bson"
)

// Number 结构体，表示带有数值的类型，支持原子操作
type Number struct {
	Value int64 `json:"number"` // 用于 BSON 序列化的字段
}

// ToBSON 将 Number 序列化为 BSON
func (num Number) ToBSON() ([]byte, error) {
	return bson.Marshal(num)
}

// Add 以原子方式增加值
func (num *Number) Add(delta int64) int64 {
	return atomic.AddInt64(&num.Value, delta)
}

// Sub 以原子方式减少值
func (num *Number) Sub(delta int64) int64 {
	return atomic.AddInt64(&num.Value, -delta)
}

// Increment 自增（+1）
func (num *Number) Increment() int64 {
	return num.Add(1)
}

// Decrement 自减（-1）
func (num *Number) Decrement() int64 {
	return num.Sub(1)
}

// Set 以原子方式设置值
func (num *Number) Set(newValue int64) {
	atomic.StoreInt64(&num.Value, newValue)
}

// Get 以原子方式获取值
func (num *Number) Get() int64 {
	return atomic.LoadInt64(&num.Value)
}

// CompareAndSwap (CAS 操作) 仅当当前值等于 old 时，才设置为 new
func (num *Number) CompareAndSwap(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&num.Value, old, new)
}
