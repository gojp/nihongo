package bloomfilter

import (
	"encoding/binary"
	"hash"
	"hash/fnv"
	"log"
	"math"
	"math/big"
)

type BloomFilter struct {
	m      uint
	k      uint
	filter *big.Int
	hasher hash.Hash64
}

// New creates a new bloom filter with m bits and k hashing functions.
func New(m uint, k uint) *BloomFilter {
	return &BloomFilter{m, k, big.NewInt(0), fnv.New64()}
}

// estimateParameters estimates the ideal parameters for a bloom filter
// with n items and fp false positive rate.
// Based on https://bitbucket.org/ww/bloom/src/829aa19d01d9/bloom.go
func estimateParameters(n uint, fp float64) (m uint, k uint) {
	m = uint(-1 * float64(n) * math.Log(fp) / math.Pow(math.Log(2), 2))
	k = uint(math.Ceil(math.Log(2) * float64(m) / float64(n)))
	return
}

// NewEstimated creates a new Bloom filter for about n items with fp
// false positive rate
func NewEstimated(n uint, fp float64) *BloomFilter {
	m, k := estimateParameters(n, fp)
	return New(m, k)
}

// baseHashes gets the two basic hash function values for data
func (f *BloomFilter) baseHashes(data []byte) (a uint32, b uint32) {
	f.hasher.Reset()
	f.hasher.Write(data)
	sum := f.hasher.Sum(nil)
	if len(sum) < 8 {
		// should never happen
		log.Println("ERROR: baseHashes found len(sum) < 8")
		return
	}
	upper := sum[0:4]
	lower := sum[4:8]
	a = binary.BigEndian.Uint32(lower)
	b = binary.BigEndian.Uint32(upper)
	return
}

// locations gets the `k` locations to set/test in the underlying bitset
func (f *BloomFilter) locations(data []byte) []int {
	locs := make([]int, f.k)
	a, b := f.baseHashes(data)
	ua := uint(a)
	ub := uint(b)
	for i := uint(0); i < f.k; i++ {
		locs[i] = int((ua + ub*i) % f.m)
	}
	return locs
}

// Add data to the Bloom Filter
func (f *BloomFilter) Add(data []byte) {
	locations := f.locations(data)
	for i := range locations {
		f.filter = f.filter.SetBit(f.filter, locations[i], 1)
	}
}

// Test checks for the presence of data in the Bloom filter
func (f *BloomFilter) Test(data []byte) bool {
	locations := f.locations(data)
	for i := uint(0); i < f.k; i++ {
		if f.filter.Bit(locations[i]) == 0 {
			return false
		}
	}
	return true
}
