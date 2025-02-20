package types

import "gopkg.in/mgo.v2/bson"

type Set struct {
	Set map[string]bool `json:"set" bson:"set"`
	TTL uint64          `json:"ttl,omitempty"`
}

// 新建一个 Set
func NewSet() *Set {
	return &Set{
		Set: make(map[string]bool),
	}
}

// 向 Set 中添加一个元素
func (s *Set) Add(value string) {
	s.Set[value] = true
}

// 检查元素是否在 Set 中
func (s *Set) Contains(value string) bool {
	return s.Set[value]
}

// 从 Set 中删除一个元素
func (s *Set) Remove(value string) {
	delete(s.Set, value)
}

// 获取 Set 中的元素数量
func (s *Set) Size() int {
	return len(s.Set)
}

// 清空 Set
func (s *Set) Clear() {
	s.TTL = 0
	s.Set = make(map[string]bool)
}

func (s *Set) ToBSON() ([]byte, error) {
	return bson.Marshal(s.Set)
}
