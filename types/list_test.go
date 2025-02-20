package types_test

import (
	"errors"
	"testing"

	"github.com/auula/wiredkv/types"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	ls := new(types.List)

	// 测试 AddItem
	ls.AddItem(10)
	ls.AddItem("hello")
	ls.AddItem(20.5)
	assert.Equal(t, 3, ls.Size(), "List 应该有 3 个元素")

	// 测试 GetItem
	item, err := ls.GetItem(1)
	assert.Nil(t, err)
	assert.Equal(t, "hello", item)

	// 测试 Remove
	err = ls.Remove("hello")
	assert.Nil(t, err)
	assert.Equal(t, 2, ls.Size(), "List 应该有 2 个元素")

	err = ls.Remove("not_exist")
	assert.Equal(t, errors.New("list item not found"), err)

	// 测试 Rnage
	ls.AddItem(30)
	rangeItems, err := ls.Rnage(0, 1)
	assert.Nil(t, err)
	assert.Equal(t, []any{10, 20.5}, rangeItems)

	// 测试 LPush
	ls.LPush(5)
	assert.Equal(t, 5, ls.List[0], "LPush 应该在列表前插入 5")

	// 测试 RPush
	ls.RPush(40)
	assert.Equal(t, 40, ls.List[ls.Size()-1], "RPush 应该在列表末尾插入 40")

	// 测试 Clear
	ls.Clear()
	assert.Equal(t, 0, ls.Size(), "Clear 之后 List 应该为空")
}
