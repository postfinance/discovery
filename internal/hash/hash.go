// Package hash implements a consistent hasher described as described in
// http://arxiv.org/pdf/1406.2294v1.pdf
package hash

import (
	"hash"
	"io"
)

// Jump is a consistent hasher.
type Jump struct {
	hasher hash.Hash64
}

// New creates a new consitent jump hasher.
func New(h hash.Hash64) *Jump {
	return &Jump{
		hasher: h,
	}
}

// Hash takes a 64 bit key and the number of buckets. It outputs a bucket
// number in the range [0, buckets).
// If the number of buckets is less than or equal to 0 then one 1 is used.
func (j *Jump) Hash(key uint64, buckets int) int {
	var b, i int64

	if buckets <= 0 {
		buckets = 1
	}

	for i < int64(buckets) {
		b = i
		key = key*2862933555777941757 + 1
		i = int64(float64(b+1) * (float64(int64(1)<<31) / float64((key>>33)+1)))
	}

	return int(b)
}

// HashString takes string as key instead of an int. It uses the configured
// hash.Hash64 hasher to create an integer from the string key.
func (j *Jump) HashString(key string, buckets int) int {
	j.hasher.Reset()

	_, err := io.WriteString(j.hasher, key)
	if err != nil {
		panic(err)
	}

	return j.Hash(j.hasher.Sum64(), buckets)
}
