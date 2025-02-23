package types

import (
	"gopkg.in/mgo.v2/bson"
)

// curl -X PUT http://192.168.31.221:2668/tables/key-03 \
//      -H "Content-Type: application/json" \
//      -H "Auth: 11111" \
//      -d '{
//        "table": {
//          "is_valid": false,
//          "items": [
//            {"id": 1, "name": "Item 1"},
//            {"id": 2, "name": "Item 2"}
//          ],
//          "meta": {
//            "version": "2.0",
//            "author": "John Doe"
//          }
//        },
//        "ttl": 10,
//      }'

// {"code":200,"message":"request processed successfully!"}

type Table struct {
	Table map[string]any `json:"table" bson:"table"`
	TTL   uint64         `json:"ttl,omitempty"`
}

// 新建一个 Table
func NewTable() *Table {
	return &Table{
		Table: make(map[string]any),
	}
}

// Clear 清空 Table 和 TTL
func (tab *Table) Clear() {
	tab.TTL = 0
	tab.Table = make(map[string]any)
}

// 向 Table 中添加一个项
func (tab *Table) AddItem(key string, value any) {
	tab.Table[key] = value
}

// 从 Table 中删除一个项
func (tab *Table) RemoveItem(key string) {
	delete(tab.Table, key)
}

// 检查 Table 中是否包含指定的键
func (tab *Table) ContainsKey(key string) bool {
	_, exists := tab.Table[key]
	return exists
}

// 从 Table 中获取一个项
func (tab *Table) GetItem(key string) any {
	if tab.ContainsKey(key) {
		return tab.Table[key]
	}
	return nil
}

// 从 Tables 查找出键为目标 key 的值，包括所有值中值
func (tab *Table) SearchItem(key string) any {
	var results []any
	if items, exists := tab.Table[key]; exists {
		results = append(results, items)
	}

	for _, item := range tab.Table {
		if innerMap, ok := item.(map[string]any); ok {
			results = append(results, searchInMap(innerMap, key)...)
		}
	}

	return results
}

func searchInMap(m map[string]any, key string) []any {
	var results []any
	if item, exists := m[key]; exists {
		results = append(results, item)
	}

	// 遍历 map，查找是否有嵌套的 map 类型
	for _, value := range m {
		if nestedMap, ok := value.(map[string]any); ok {
			// 递归查找嵌套的 map
			results = append(results, searchInMap(nestedMap, key)...)
		}
	}

	return results
}

// 获取 Table 中的元素个数
func (tab *Table) Size() int {
	return len(tab.Table)
}

func (tab Table) ToBSON() ([]byte, error) {
	return bson.Marshal(tab)
}
