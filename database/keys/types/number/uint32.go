package number

import (
	"encoding/binary"
	"io"
)

// Uint32 is a codec for encoding and decoding 32-bit unsigned integers.
type Uint32 struct {
	value uint32
}

// Set sets the value as a uint32.
func (c *Uint32) Set(value uint32) {
	c.value = value
}

// Get gets the value as a uint32.
func (c *Uint32) Get() uint32 {
	return c.value
}

// SetInt sets the value as an int, converting it to uint32.
// Values outside the range of uint32 (0–4294967295) will be truncated.
func (c *Uint32) SetInt(value int) {
	c.value = uint32(value)
}

// Int gets the value as an int, converted from uint32.
func (c *Uint32) Int() int {
	return int(c.value)
}

// MarshalWrite writes the uint32 value to the provided writer in BigEndian order.
func (c *Uint32) MarshalWrite(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, c.value)
}

// UnmarshalRead reads a uint32 value from the provided reader in BigEndian order.
func (c *Uint32) UnmarshalRead(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &c.value)
}

type Uint32s []*Uint32

// Union computes the union of the current Uint32s slice with another Uint32s slice. The result
// contains all unique elements from both slices.
func (s Uint32s) Union(other Uint32s) Uint32s {
	valueMap := make(map[uint32]bool)
	var result Uint32s

	// Add elements from the current Uint32s slice to the result
	for _, item := range s {
		val := item.Get()
		if !valueMap[val] {
			valueMap[val] = true
			result = append(result, item)
		}
	}

	// Add elements from the other Uint32s slice to the result
	for _, item := range other {
		val := item.Get()
		if !valueMap[val] {
			valueMap[val] = true
			result = append(result, item)
		}
	}

	return result
}

// Intersection computes the intersection of the current Uint32s slice with another Uint32s
// slice. The result contains only the elements that exist in both slices.
func (s Uint32s) Intersection(other Uint32s) Uint32s {
	valueMap := make(map[uint32]bool)
	var result Uint32s

	// Add all elements from the other Uint32s slice to the map
	for _, item := range other {
		valueMap[item.Get()] = true
	}

	// Check for common elements in the current Uint32s slice
	for _, item := range s {
		val := item.Get()
		if valueMap[val] {
			result = append(result, item)
		}
	}

	return result
}

// Difference computes the difference of the current Uint32s slice with another Uint32s slice.
// The result contains only the elements that are in the current slice but not in the other
// slice.
func (s Uint32s) Difference(other Uint32s) Uint32s {
	valueMap := make(map[uint32]bool)
	var result Uint32s

	// Mark all elements in the other Uint32s slice
	for _, item := range other {
		valueMap[item.Get()] = true
	}

	// Add elements from the current Uint32s slice that are not in the other Uint32s slice
	for _, item := range s {
		val := item.Get()
		if !valueMap[val] {
			result = append(result, item)
		}
	}

	return result
}
