package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZSet_Add(t *testing.T) {
	zset := NewZSet()

	// Test adding new elements
	zset.Add("item1", 10.5)
	assert.Equal(t, 1, zset.Size())
	assert.Contains(t, zset.ZSet, "item1")
	assert.Equal(t, 10.5, zset.ZSet["item1"])

	// Test adding another element with a different score
	zset.Add("item2", 20.5)
	assert.Equal(t, 2, zset.Size())
	assert.Contains(t, zset.ZSet, "item2")
	assert.Equal(t, 20.5, zset.ZSet["item2"])

	// Test updating score of an existing element
	zset.Add("item1", 15.0)
	assert.Equal(t, 2, zset.Size())
	assert.Equal(t, 15.0, zset.ZSet["item1"])
}

func TestZSet_Remove(t *testing.T) {
	zset := NewZSet()
	zset.Add("item1", 10.5)

	// Test removing an element
	zset.Remove("item1")
	assert.Equal(t, 0, zset.Size())
	assert.NotContains(t, zset.ZSet, "item1")

	// Test removing a non-existing element
	zset.Remove("item2")
	assert.Equal(t, 0, zset.Size())
}

func TestZSet_Get(t *testing.T) {
	zset := NewZSet()
	zset.Add("item1", 10.5)

	// Test getting the score of an existing element
	score, exists := zset.Get("item1")
	assert.True(t, exists)
	assert.Equal(t, 10.5, score)

	// Test getting the score of a non-existing element
	_, exists = zset.Get("item2")
	assert.False(t, exists)
}

func TestZSet_GetRank(t *testing.T) {
	zset := NewZSet()
	zset.Add("item1", 10.5)
	zset.Add("item2", 20.5)

	// Test getting rank of an element
	rank, exists := zset.GetRank("item1")
	assert.True(t, exists)
	assert.Equal(t, 1, rank)

	// Test getting rank of a non-existing element
	_, exists = zset.GetRank("item3")
	assert.False(t, exists)
}

func TestZSet_GetRange(t *testing.T) {
	zset := NewZSet()
	zset.Add("item1", 10.5)
	zset.Add("item2", 20.5)
	zset.Add("item3", 15.0)

	// Test getting elements in a score range
	rangeItems := zset.GetRange(10.0, 20.0)
	assert.Equal(t, []string{"item3", "item1"}, rangeItems)

	// Test with a score range that doesn't match any element
	rangeItems = zset.GetRange(30.0, 40.0)
	assert.Empty(t, rangeItems)
}

func TestZSet_Sort(t *testing.T) {
	zset := NewZSet()
	zset.Add("item1", 10.5)
	zset.Add("item2", 20.5)
	zset.Add("item3", 15.0)

	// Test if sortedScores reflects correct order
	assert.Equal(t, []string{"item2", "item3", "item1"}, zset.sortedScores)
}

func TestZSet_Clear(t *testing.T) {
	zset := NewZSet()
	zset.Add("item1", 10.5)
	zset.Clear()

	// Test clear functionality
	assert.Equal(t, 0, zset.Size())
	assert.Equal(t, uint64(0), zset.TTL)
	assert.Empty(t, zset.ZSet)
}

func TestZSet_ToBSON(t *testing.T) {
	zset := NewZSet()
	zset.Add("item1", 10.5)

	// Test ToBSON
	_, err := zset.ToBSON()
	assert.NoError(t, err)
}
