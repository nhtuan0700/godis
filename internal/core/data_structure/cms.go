package data_structure

import (
	"math"

	"github.com/spaolacci/murmur3"
)

// Log10PointFive is precoumputed value for log_10(0.5).
const Log10PointFive = -0.30102999566

// CMS is the Count-Min Sketch data structure.
// The counter field has been changed to a 2D slice for better clarity and indexing.
type CMS struct {
	width uint32
	depth uint32
	// counter is now 2D slice of uint64. The outer slice represents the rows (depth),
	// and the inner slice represents the columns (width).
	counter [][]uint64
}

// CreateCMS initializes a new Count-Min Sketch with a give width and depth.
func CreateCMS(w, d uint32) *CMS {
	cms := &CMS{
		width: w,
		depth: d,
	}

	// Initialize 2D slice.
	// We create a slice of slices, where outer slice has 'd' elements (for depth).
	counter := make([][]uint64, d)
	// Then we loop through each "row" and initalize a slice of 'w' elements for the width.
	for i := uint32(0); i < d; i++ {
		counter[i] = make([]uint64, w)
	}
	cms.counter = counter

	return cms
}

// CalcCMSDim calculates the dimensions (width and depth) of the CMS
// based on the desired error rate and probability.
func CalcCMSDim(errRate float64, probability float64) (uint32, uint32) {
	w := uint32(math.Ceil(2.0 / errRate))
	d := uint32(math.Ceil(math.Log10(probability) / Log10PointFive))
	return w, d
}

// calcHash calculate 32-bit hash for the given item and seed.
func (c *CMS) calcHash(item string, seed uint32) uint32 {
	hasher := murmur3.New32WithSeed(seed)
	hasher.Write([]byte(item))
	return hasher.Sum32()
}


// IncrBy increments the count for an item by the specific value.
// It return estimated count for the item after the increment.
func (c *CMS) IncrBy(item string, value uint64) uint64 {
	var minCount uint64 = math.MaxUint64

	// Loop through each row of the 2D array.
	for i := uint32(0); i < c.depth; i++ {
		// Calculate a new hash for each row using the row index as the seed.
		hash := c.calcHash(item, i)
		// Use the has to get the column index within a row
		j := hash % c.width

		// Safely add the value to prevent the overflow.
		if math.MaxUint64-c.counter[i][j] < value {
			c.counter[i][j] = math.MaxUint64
		} else {
			c.counter[i][j] += value
		}

		// Keep track of the minimum count across all rows.
		if c.counter[i][j] < minCount {
			minCount = c.counter[i][j]
		}
	}

	return minCount
}

// Count return the estimated ccount for an item.
// It retrieves the minimum count across all hash functions to provide the most accurate estimate.
func (c *CMS) Count(item string) uint64 {
	var minCount uint64 = math.MaxUint64

	// Loop throught each row of the 2D array.
	for i := uint32(0); i < c.depth; i++ {
		// Calculate the hash for this row.
		hash := c.calcHash(item, i)
		// Determine the column index.
		j := hash % c.width

		// Find the minimum count across all rows.
		if c.counter[i][j] < minCount {
			minCount = c.counter[i][j]
		}
	}

	return minCount
}
