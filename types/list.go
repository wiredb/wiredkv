package types

import (
	"errors"

	"gopkg.in/mgo.v2/bson"
)

type List struct {
	List []any  `json:"list" bson:"list" binding:"required"`
	TTL  uint64 `json:"ttl,omitempty"`
}

func NewList() *List {
	return new(List)
}

// AddItem 向 List 中添加新项目
func (ls *List) AddItem(item any) {
	ls.List = append(ls.List, item)
}

// Remove 从 List 中删除指定的项目
func (ls *List) Remove(item any) error {
	for i, v := range ls.List {
		if v == item {
			ls.List = append(ls.List[:i], ls.List[i+1:]...)
			return nil
		}
	}
	return errors.New("list item not found")
}

// GetItem 获取 List 中指定索引的项目
func (ls *List) GetItem(index int) (any, error) {
	if index < 0 || index >= len(ls.List) {
		return nil, errors.New("list index out of bounds")
	}
	return ls.List[index], nil
}

func (ls *List) Rnage(statIndex, endIndex int) ([]any, error) {
	var result []any
	for i, v := range ls.List {
		if i >= statIndex && i <= endIndex {
			result = append(result, v)
		}
	}
	return result, nil
}

func (ls *List) LPush(item any) {
	ls.List = append([]any{item}, ls.List...)
}

func (ls *List) RPush(item any) {
	ls.List = append(ls.List, item)
}

func (ls *List) Size() int {
	return len(ls.List)
}

func (ls *List) Clear() {
	ls.TTL = 0
	ls.List = make([]any, 0)
}

func (ls List) ToBSON() ([]byte, error) {
	return bson.Marshal(ls)
}
