package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestTextOperations(t *testing.T) {
	text := NewText("Hello, World!")
	assert.Equal(t, 13, text.Size())
	assert.True(t, text.Contains("World"))
	assert.False(t, text.Contains("Leon Ding"))

	text.Append(" and more")
	assert.Equal(t, "Hello, World! and more", text.Content)

	clone := text.Clone()
	assert.Equal(t, text.Content, clone.Content)
	assert.Equal(t, text.TTL, clone.TTL)

	text.Clear()
	assert.Equal(t, "", text.Content)
	assert.Equal(t, uint64(0), text.TTL)

	bsonData, err := text.ToBSON()
	assert.NoError(t, err)

	var decodedText Text
	err = bson.Unmarshal(bsonData, &decodedText)
	assert.NoError(t, err)
	assert.Equal(t, text, &decodedText)

	other := NewText("other text content.")
	assert.False(t, false, text, other)
}
