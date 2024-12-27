package vfs

import (
	"errors"
	"fmt"

	"github.com/auula/vasedb/types"
)

type Kind int8

const (
	Set Kind = iota
	ZSet
	List
	Text
	Tables
	Binary
	Number
	Unknown
)

type Segment struct {
	kind Kind
	data []byte
}

type Serializable interface {
	ToBytes() []byte
}

// NewSegment 使用数据类型初始化并返回对应的 Segment
func NewSegment(data Serializable) (*Segment, error) {
	kind, err := toKind(data)
	if err != nil {
		return nil, fmt.Errorf("unsupported data type: %w", err)
	}

	// 如果类型不匹配，则返回错误
	return &Segment{
		kind: kind,
		data: data.ToBytes(),
	}, nil
}

func (s *Segment) Kind() Kind {
	return s.kind
}

func (s *Segment) Size() int {
	return len(s.data)
}

func (s *Segment) ToBytes() []byte {

	return []byte{}
}

func (s *Segment) ToSet() *types.Set {
	if s.kind != Set {
		return nil
	}
	// 假设您的数据是 JSON 或某种结构体，可以进行反序列化
	var set types.Set
	// Deserialize s.data to set (具体根据类型定义来做反序列化)
	return &set
}

func (s *Segment) ToZSet() *types.ZSet {
	return nil
}

func (s *Segment) ToText() *types.Text {
	return nil
}

func (s *Segment) ToList() *types.List {
	return nil
}

func (s *Segment) ToTables() *types.Tables {
	return nil
}

func (s *Segment) ToBinary() *types.Binary {
	return nil
}

func (s *Segment) ToNumber() *types.Number {
	return nil
}

func (s *Segment) TTL() {

}

// 将类型映射为 Kind 的辅助函数
func toKind(data Serializable) (Kind, error) {
	switch data.(type) {
	case *types.Set:
		return Set, nil
	case *types.ZSet:
		return ZSet, nil
	case *types.List:
		return List, nil
	case *types.Text:
		return Text, nil
	case *types.Tables:
		return Tables, nil
	case *types.Binary:
		return Binary, nil
	case *types.Number:
		return Number, nil
	default:
		return Unknown, errors.New("unknown data type")
	}
}
