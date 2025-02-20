package vfs

import (
	"testing"
	"time"

	"github.com/auula/wiredkv/types"
	"github.com/stretchr/testify/assert"
)

func TestNewSegment(t *testing.T) {
	// Test valid Set type
	set := types.Set{
		Set: map[string]bool{
			"item1": true,
			"item2": true,
		},
	}

	// Create a new segment for the Set type
	segment, err := NewSegment("mock-key", &set, 1000)
	assert.NoError(t, err)                                    // Ensure no error
	assert.NotNil(t, segment)                                 // Ensure segment is created
	assert.Equal(t, "mock-key", string(segment.Key))          // Ensure the key is set correctly
	assert.Equal(t, uint32(len("mock-key")), segment.KeySize) // Ensure the key size is correct
	assert.Equal(t, uint32(21), segment.ValueSize)            // Ensure the value size is correct
}

func TestNewTombstoneSegment(t *testing.T) {
	// Create a Tombstone segment
	segment := NewTombstoneSegment("mock-key")

	// Ensure the segment is of Tombstone type and has expected fields
	assert.Equal(t, Unknown, segment.Type)                    // Tombstone should have Unknown type
	assert.Equal(t, int8(1), segment.Tombstone)               // Tombstone should be marked as 1
	assert.Equal(t, "mock-key", string(segment.Key))          // Ensure the key is set correctly
	assert.Equal(t, uint32(len("mock-key")), segment.KeySize) // Ensure the key size is correct
}

func TestSegmentSize(t *testing.T) {
	// Create a Set type data for testing
	set := types.Set{
		Set: map[string]bool{
			"item1": true,
			"item2": true,
		},
	}

	// Create a segment for the Set type
	segment, err := NewSegment("mock-key", &set, 1000)
	assert.NoError(t, err)

	// Ensure the size is calculated correctly
	assert.Equal(t, uint32(59), segment.Size())
}

func TestToSet(t *testing.T) {
	// Create a Set type Segment
	setData := types.Set{
		Set: map[string]bool{
			"item1": true,
			"item2": true,
		},
	}
	segment, err := NewSegment("mock-key", &setData, 1000)
	assert.NoError(t, err)

	// Convert the segment to Set
	set, err := segment.ToSet()
	assert.NoError(t, err)                // Ensure no error
	assert.Equal(t, setData.Set, set.Set) // Ensure the Set values match
}

func TestTTL(t *testing.T) {
	// Create a Segment with TTL
	set := types.Set{
		Set: map[string]bool{
			"item1": true,
			"item2": true,
		},
	}
	segment, err := NewSegment("mock-key", &set, 1) // TTL = 1 second
	assert.NoError(t, err)

	// Wait 1 second
	time.Sleep(time.Second)

	// Test TTL, it should return a value close to 0
	ttl := segment.TTL()
	assert.True(t, ttl <= 0) // Ensure TTL is <= 0 after expiration
}
