package data_structure

import (
	"math"

	"github.com/spaolacci/murmur3"
)

const Ln2 float64 = 0.693147180559945
const Ln2Square float64 = 0.480453013918201 // ln(2)^2
const ABigSeed uint32 = 0x9747b28c          // 2,538,058,380

type BloomFilter struct {
	Hashes      int
	Entries     uint64
	Error       float64
	bitPerEntry float64
	bf          []uint8
	bits        uint64
	bytes       uint64
}

type HashValue struct {
	a uint64
	b uint64
}

func calcBpe(errorRate float64) float64 {
	num := math.Log(errorRate)
	return math.Abs(-num / Ln2Square)
}

/*
http://en.wikipedia.org/wiki/Bloom_filter
- Optimal number of bits is: bits = (entries * ln(error)) / ln(2)^2
- bitPerEntry = bits/entries
- Optimal number of hash functions is: hashes = bitPerEntry * ln(2)
*/
func CreateBloomFilter(entries uint64, errorRate float64) *BloomFilter {
	bloom := &BloomFilter{
		Entries: entries,
		Error:   errorRate,
	}

	bloom.bitPerEntry = calcBpe(errorRate)
	bits := uint64(bloom.bitPerEntry * float64(entries))
	if bits%64 != 0 {
		bloom.bytes = ((bits / 64) + 1) * 8
	} else {
		bloom.bytes = bits / 8
	}

	bloom.bits = bloom.bytes * 8
	bloom.Hashes = int(math.Ceil(Ln2 * bloom.bitPerEntry))
	bloom.bf = make([]uint8, bloom.bytes)

	return bloom
}

func (b *BloomFilter) CalcHash(entry string) HashValue {
	hasher := murmur3.New128WithSeed(ABigSeed)
	hasher.Write([]byte(entry))
	x, y := hasher.Sum128()
	return HashValue{
		a: x,
		b: y,
	}
}

func (b *BloomFilter) Add(entry string) bool {
	initHash := b.CalcHash(entry)
	return b.AddHash(initHash)
}

func (b *BloomFilter) Exist(entry string) bool {
	initHash := b.CalcHash(entry)
	return b.ExistHash(initHash)
}

func (b *BloomFilter) AddHash(hashValue HashValue) bool {
	foundUnset := false
	
	for i := 0; i < b.Hashes; i++ {
		hash := (hashValue.a + hashValue.b*uint64(i)) % b.bits
		bytePos := hash >> 3    // div 8
		mask := 1 << (hash % 8) // mask at candidate byte position
		if b.bf[bytePos]&uint8(mask) == 0 {
			b.bf[bytePos] |= uint8(mask)
			foundUnset = true
		}
	}

	return foundUnset
}

func (b *BloomFilter) ExistHash(hashValue HashValue) bool {
	for i := 0; i < b.Hashes; i++ {
		hash := (hashValue.a + hashValue.b*uint64(i)) % b.bits
		bytePos := hash >> 3    // div 8
		mask := 1 << (hash % 8) // mask at candidate byte position
		if b.bf[bytePos]&uint8(mask) == 0 {
			return false
		}
	}

	return true
}
