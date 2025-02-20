package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestNumberOperations(t *testing.T) {
	num := NewNumber(10)
	assert.Equal(t, int64(10), num.Get())

	num.Add(5)
	assert.Equal(t, int64(15), num.Get())

	num.Sub(3)
	assert.Equal(t, int64(12), num.Get())

	num.Increment()
	assert.Equal(t, int64(13), num.Get())

	num.Decrement()
	assert.Equal(t, int64(12), num.Get())

	num.Set(100)
	assert.Equal(t, int64(100), num.Get())

	success := num.CompareAndSwap(100, 200)
	assert.True(t, success)
	assert.Equal(t, int64(200), num.Get())

	success = num.CompareAndSwap(100, 300)
	assert.False(t, success)
	assert.Equal(t, int64(200), num.Get())

	bsonData, err := num.ToBSON()
	assert.NoError(t, err)

	var decodedNumber Number
	err = bson.Unmarshal(bsonData, &decodedNumber)
	assert.NoError(t, err)
	assert.Equal(t, num.Get(), decodedNumber.Get())
}
