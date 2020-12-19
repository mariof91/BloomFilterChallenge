package data

import (
	"github.com/spaolacci/murmur3"
	"hash"
	"hash/fnv"
	"math"
)

type Filter struct{
	bitset []bool      	//bloom-filter bitset
	k      uint        	// Number of hash values
	n      uint         // Number of elements in the filter
	m      uint         // Size of the bloom filter bitset
	hashFuncs[]hash.Hash64           // The hash functions
}

//type Stat struct {
//	size uint 							`json:size`
//	functions uint						`json:functions`
//	count uint							`json:count`
//	falsePositiveProbability float64	`json:falsePositiveProbability`
//}

type Stat struct {
	Size uint 							`json:"size"`
	Functions uint						`json:"functions"`
	Count uint							`json:"count"`
	FalsePositiveProbability float64	`json:"falsePositiveProbability"`
}


// Returns a new BloomFilter object,
func New(size uint) *Filter {
	return &Filter{
		bitset: make([]bool, size),
		k: 3,
		m: size,
		n: uint(0),
		hashFuncs:[]hash.Hash64{murmur3.New64(),fnv.New64(),fnv.New64a()},
	}
}

type FilterInterface interface {
	Add(item []byte)     // Adds the item into the Set
	Test(item []byte) bool  // Check if items is maybe in the Set
}

// Adds the item into the bloom filter set by hashing in over the . // hash functions
func (f *Filter) Add(item []byte) {
	hashes := f.hashValues(item)
	for i:=uint(0); i < f.k; i++ {
		position := uint(hashes[i]) % f.m
		f.bitset[position] = true
		}
	f.n++
}

// Calculates all the hash values by applying in the item over the // hash functions
func (f *Filter) hashValues(item []byte) []uint64  {
	var result []uint64
	for _, hashFunc := range f.hashFuncs {
		hashFunc.Write(item)
		result = append(result, hashFunc.Sum64())
		hashFunc.Reset()
	}
	return result
}

// Test if the item into the bloom filter is set by hashing in over // the hash functions
func (f *Filter) Test(item []byte)  bool{
	hashes := f.hashValues(item)
	for i:=uint(0); i < f.k; i++ {
		position := uint(hashes[i]) % f.m
		if !f.bitset[uint(position)] {
			return false
		}
	}
	return true
}

func (f Filter) GetStat() *Stat {

	return &Stat{
		Size:                     f.m,
		Functions:                f.k,
		Count:                    f.n,
		FalsePositiveProbability: f.CalculateFalsePositive(),
	}
}

func (f Filter) CalculateFalsePositive()float64{
	m:= float64(f.m)
	n:= float64(f.n)
	k:= float64(f.k)

	//pow(1 - exp(-k / (m / n)), k)
	return math.Pow(1 - math.Exp(-k / (m / n)), k)
}