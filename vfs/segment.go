package vfs

import (
	"errors"
	"fmt"
	"time"

	"github.com/auula/wiredkv/types"
	"gopkg.in/mgo.v2/bson"
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

// | DEL 1 | KIND 1 | EAT 8 | CAT 8 | KLEN 8 | VLEN 8 | KEY ? | VALUE ? | CRC32 4 |
type Segment struct {
	Tombstone int8
	Type      Kind
	ExpiredAt uint64
	CreatedAt uint64
	KeySize   uint32
	ValueSize uint32
	Key       []byte
	Value     []byte
}

type Serializable interface {
	ToBSON() ([]byte, error)
}

// NewSegment 使用数据类型初始化并返回对应的 Segment
func NewSegment(key string, data Serializable, ttl uint64) (*Segment, error) {
	kind, err := toKind(data)
	if err != nil {
		return nil, fmt.Errorf("unsupported data type: %w", err)
	}

	timestamp, expiredAt := uint64(time.Now().Unix()), uint64(0)
	if ttl > 0 {
		expiredAt = uint64(time.Now().Add(time.Second * time.Duration(ttl)).Unix())
	}

	bytes, err := data.ToBSON()
	if err != nil {
		return nil, err
	}

	// 这个是通过 transformer 编码之后的
	encodedata, err := transformer.Encode(bytes)
	if err != nil {
		return nil, fmt.Errorf("transformer encode: %w", err)
	}

	// 如果类型不匹配，则返回错误
	return &Segment{
		Type:      kind,
		Tombstone: 0,
		CreatedAt: timestamp,
		ExpiredAt: expiredAt,
		KeySize:   uint32(len(key)),
		ValueSize: uint32(len(encodedata)),
		Key:       []byte(key),
		Value:     encodedata,
	}, nil

}

func NewTombstoneSegment(key string) *Segment {
	timestamp, expiredAt := uint64(time.Now().Unix()), uint64(0)
	return &Segment{
		Type:      Unknown,
		Tombstone: 1,
		CreatedAt: timestamp,
		ExpiredAt: expiredAt,
		KeySize:   uint32(len(key)),
		ValueSize: 0,
		Key:       []byte(key),
		Value:     []byte{},
	}
}

func (s *Segment) IsTombstone() bool {
	return s.Tombstone == 1
}

func (s *Segment) Size() uint32 {
	// 计算一整块记录的大小，+4 CRC 校验码占用 4 个字节
	return 26 + s.KeySize + s.ValueSize + 4
}

func (s *Segment) ToSet() (*types.Set, error) {
	if s.Type != Set {
		return nil, fmt.Errorf("not support conversion to set type")
	}
	var set types.Set
	err := bson.Unmarshal(s.Value, &set)
	if err != nil {
		return nil, err
	}
	return &set, nil
}

func (s *Segment) ToZSet() (*types.ZSet, error) {
	if s.Type != ZSet {
		return nil, fmt.Errorf("not support conversion to zset type")
	}
	var zset types.ZSet
	err := bson.Unmarshal(s.Value, &zset)
	if err != nil {
		return nil, err
	}
	return &zset, nil
}

func (s *Segment) ToText() (*types.Text, error) {
	if s.Type != Text {
		return nil, fmt.Errorf("not support conversion to text type")
	}
	var text types.Text
	err := bson.Unmarshal(s.Value, &text)
	if err != nil {
		return nil, err
	}
	return &text, nil
}

func (s *Segment) ToList() (*types.List, error) {
	if s.Type != List {
		return nil, fmt.Errorf("not support conversion to list type")
	}
	var list types.List
	err := bson.Unmarshal(s.Value, &list)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func (s *Segment) ToTables() (*types.Tables, error) {
	if s.Type != Tables {
		return nil, fmt.Errorf("not support conversion to tables type")
	}
	var tables types.Tables
	err := bson.Unmarshal(s.Value, &tables)
	if err != nil {
		return nil, err
	}
	return &tables, nil
}

func (s *Segment) ToBinary() (*types.Binary, error) {
	if s.Type != Binary {
		return nil, fmt.Errorf("not support conversion to binary type")
	}
	var bin types.Binary
	err := bson.Unmarshal(s.Value, &bin)
	if err != nil {
		return nil, err
	}
	return &bin, nil
}

func (s *Segment) ToNumber() (*types.Number, error) {
	if s.Type != Number {
		return nil, fmt.Errorf("not support conversion to number type")
	}
	var number types.Number
	err := bson.Unmarshal(s.Value, &number)
	if err != nil {
		return nil, err
	}
	return &number, nil
}

func (s *Segment) TTL() int64 {
	now := uint64(time.Now().Unix())
	if s.ExpiredAt > 0 && s.ExpiredAt > now {
		return int64(s.ExpiredAt - now)
	}
	return -1
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
	case types.Tables:
		return Tables, nil
	case *types.Binary:
		return Binary, nil
	case *types.Number:
		return Number, nil
	default:
		return Unknown, errors.New("unknown data type")
	}
}
