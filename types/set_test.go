package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestSetOperations(t *testing.T) {
	set := NewSet()
	assert.NotNil(t, set)
	assert.Equal(t, 0, set.Size())

	set.Add("item1")
	assert.True(t, set.Contains("item1"))

	assert.False(t, set.Contains("item2"))
	set.Add("item2")
	assert.True(t, set.Contains("item2"))

	set.Remove("item1")
	assert.False(t, set.Contains("item1"))
	assert.Equal(t, 1, set.Size())

	set.Clear()
	assert.Equal(t, 0, set.Size())
	assert.False(t, set.Contains("item1"))
	assert.False(t, set.Contains("item2"))

	set.Add("item1")
	set.Add("item2")

	bsonData, err := set.ToBSON()
	assert.NoError(t, err)

	var decodedMap map[string]bool
	err = bson.Unmarshal(bsonData, &decodedMap)
	assert.NoError(t, err)
	assert.Equal(t, set.Set, decodedMap)
}
