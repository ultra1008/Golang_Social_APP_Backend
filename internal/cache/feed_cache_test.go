package cache

import (
	"testing"

	"github.com/niklod/highload-social-network/internal/user/post"
	"github.com/stretchr/testify/assert"
)

func TestFeedCache_WriteRead(t *testing.T) {
	cache := NewFeedCache()

	testKey, testValue := 1, post.Post{}

	cache.Write(testKey, testValue)

	v, ok := cache.Read(testKey)

	assert.Equal(t, true, ok)
	assert.Equal(t, v, testValue)
}

func TestFeedCache_WriteRead_ReplacePrevValue(t *testing.T) {
	cache := NewFeedCache()

	testKey := 1
	testValue := post.Feed{}
	testValue2 := post.Feed{{ID: 1}}

	cache.Write(testKey, testValue)

	v, ok := cache.Read(testKey)

	assert.NotEqual(t, false, ok)
	assert.Equal(t, v, testValue)

	c := v.(post.Feed)
	assert.Equal(t, len(c), len(testValue))

	cache.Write(testKey, testValue2)

	v, ok = cache.Read(testKey)

	c = v.(post.Feed)
	assert.NotEqual(t, false, ok)
	assert.Equal(t, c, testValue2)
	assert.Equal(t, len(c), len(testValue2))
}

func TestFeedCache_Read_EmptyResult(t *testing.T) {
	cache := NewFeedCache()

	testKey := 1

	v, ok := cache.Read(testKey)

	assert.Equal(t, false, ok)
	assert.Nil(t, v)
}
