package hash

import (
	"fmt"
	"hash/crc64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	var tt = []struct {
		key      uint64
		buckets  int
		expected int
	}{
		{1, 1, 0},
		{42, 57, 43},
		{0xDEAD10CC, 1, 0},
		{0xDEAD10CC, 666, 361},
		{256, 1024, 520},
		{0, -10, 0},
		{0xDEAD10CC, -666, 0},
	}

	h := New(nil)

	for i := range tt {
		tc := tt[i]
		t.Run(fmt.Sprintf("key=%d, buckets=%d", tc.key, tc.buckets), func(t *testing.T) {
			assert.Equal(t, h.Hash(tc.key, tc.buckets), tc.expected)
		})
	}
}

func TestHashString(t *testing.T) {
	var tt = []struct {
		key      string
		buckets  int
		expected int
	}{
		{"localhost", 10, 6},
		{"ёлка", 3, 0},
		{"ветер", 10, 7},
	}

	h := New(crc64.New(crc64.MakeTable(0xC96C5795D7870F42)))

	for i := range tt {
		tc := tt[i]
		t.Run(fmt.Sprintf("key=%s, buckets=%d", tc.key, tc.buckets), func(t *testing.T) {
			assert.Equal(t, h.HashString(tc.key, tc.buckets), tc.expected)
		})
	}
}
