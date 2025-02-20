package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZSet(t *testing.T) {
	zset := NewZSet()

	// 测试 Add 方法
	zset.Add("a", 10)
	zset.Add("b", 20)
	zset.Add("c", 15)
	assert.Equal(t, 3, zset.Size(), "ZSet 应该有 3 个元素")

	// 测试 Get 方法
	score, exists := zset.Get("b")
	assert.True(t, exists, "元素 'b' 应该存在")
	assert.Equal(t, 20.0, score, "元素 'b' 的分数应该是 20")

	// 测试 GetRank 方法
	rank, found := zset.GetRank("c")
	assert.True(t, found, "元素 'c' 应该有排名")
	assert.Equal(t, 1, rank, "元素 'c' 在排名中应该是第 1 (从 0 开始)")

	// 测试 Remove 方法
	zset.Remove("b")
	_, exists = zset.Get("b")
	assert.False(t, exists, "元素 'b' 应该被删除")
	assert.Equal(t, 2, zset.Size(), "ZSet 现在应该有 2 个元素")

	// 测试 GetRange 方法
	result := zset.GetRange(10, 15)
	assert.ElementsMatch(t, []string{"a", "c"}, result, "GetRange 应该返回正确的元素")

	// 测试 Clear 方法
	zset.Clear()
	assert.Equal(t, 0, zset.Size(), "Clear 之后 ZSet 应该为空")
}
