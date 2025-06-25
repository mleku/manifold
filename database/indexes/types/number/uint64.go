package number

import (
	"encoding/binary"
	"io"
)

// Uint64 is a codec for encoding and decoding 64-bit unsigned integers.
type Uint64 struct {
	value uint64
}

// Set sets the value as a uint64.
func (c *Uint64) Set(value uint64) {
	c.value = value
}

// Get gets the value as a uint64.
func (c *Uint64) Get() uint64 {
	return c.value
}

// SetInt sets the value as an int, converting it to uint64.
// Values outside the range of uint64 are truncated.
func (c *Uint64) SetInt(value int) {
	c.value = uint64(value)
}

// Int gets the value as an int, converted from uint64. May truncate if the value exceeds the
// range of int.
func (c *Uint64) Int() int {
	return int(c.value)
}

// MarshalWrite writes the uint64 value to the provided writer in BigEndian order.
func (c *Uint64) MarshalWrite(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, c.value)
}

// UnmarshalRead reads a uint64 value from the provided reader in BigEndian order.
func (c *Uint64) UnmarshalRead(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &c.value)
}

type Uint64s []*Uint64

// Union computes the union of the current Uint64s slice with another Uint64s slice. The result
// contains all unique elements from both slices.
func (s Uint64s) Union(other Uint64s) Uint64s {
	valueMap := make(map[uint64]bool)
	var result Uint64s

	// Add elements from the current Uint64s slice to the result
	for _, item := range s {
		val := item.Get()
		if !valueMap[val] {
			valueMap[val] = true
			result = append(result, item)
		}
	}

	// Add elements from the other Uint64s slice to the result
	for _, item := range other {
		val := item.Get()
		if !valueMap[val] {
			valueMap[val] = true
			result = append(result, item)
		}
	}

	return result
}

// Intersection computes the intersection of the current Uint64s slice with another Uint64s
// slice. The result contains only the elements that exist in both slices.
func (s Uint64s) Intersection(other Uint64s) Uint64s {
	valueMap := make(map[uint64]bool)
	var result Uint64s

	// Add all elements from the other Uint64s slice to the map
	for _, item := range other {
		valueMap[item.Get()] = true
	}

	// Check for common elements in the current Uint64s slice
	for _, item := range s {
		val := item.Get()
		if valueMap[val] {
			result = append(result, item)
		}
	}

	return result
}

// Difference computes the difference of the current Uint64s slice with another Uint64s slice.
// The result contains only the elements that are in the current slice but not in the other
// slice.
func (s Uint64s) Difference(other Uint64s) Uint64s {
	valueMap := make(map[uint64]bool)
	var result Uint64s

	// Mark all elements in the other Uint64s slice
	for _, item := range other {
		valueMap[item.Get()] = true
	}

	// Add elements from the current Uint64s slice that are not in the other Uint64s slice
	for _, item := range s {
		val := item.Get()
		if !valueMap[val] {
			result = append(result, item)
		}
	}

	return result
}
