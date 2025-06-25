package number

import (
	"encoding/binary"
	"io"
)

// Uint16 is a codec for encoding and decoding 16-bit unsigned integers.
type Uint16 struct {
	value uint16
}

// Set sets the value as a uint16.
func (c *Uint16) Set(value uint16) {
	c.value = value
}

// Get gets the value as a uint16.
func (c *Uint16) Get() uint16 {
	return c.value
}

// SetInt sets the value as an int, converting it to uint16. Truncates values outside uint16 range (0-65535).
func (c *Uint16) SetInt(value int) {
	c.value = uint16(value)
}

// GetInt gets the value as an int, converted from uint16.
func (c *Uint16) GetInt() int {
	return int(c.value)
}

// MarshalWrite writes the uint16 value to the provided writer in BigEndian order.
func (c *Uint16) MarshalWrite(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, c.value)
}

// UnmarshalRead reads a uint16 value from the provided reader in BigEndian order.
func (c *Uint16) UnmarshalRead(r io.Reader) error {
	return binary.Read(r, binary.BigEndian, &c.value)
}

type Uint16s []*Uint16

// Union computes the union of the current Uint16s slice with another Uint16s slice. The result
// contains all unique elements from both slices.
func (s Uint16s) Union(other Uint16s) Uint16s {
	valueMap := make(map[uint16]bool)
	var result Uint16s

	// Add elements from the current Uint16s slice to the result
	for _, item := range s {
		val := item.Get()
		if !valueMap[val] {
			valueMap[val] = true
			result = append(result, item)
		}
	}

	// Add elements from the other Uint16s slice to the result
	for _, item := range other {
		val := item.Get()
		if !valueMap[val] {
			valueMap[val] = true
			result = append(result, item)
		}
	}

	return result
}

// Intersection computes the intersection of the current Uint16s slice with another Uint16s
// slice. The result contains only the elements that exist in both slices.
func (s Uint16s) Intersection(other Uint16s) Uint16s {
	valueMap := make(map[uint16]bool)
	var result Uint16s

	// Add all elements from the other Uint16s slice to the map
	for _, item := range other {
		valueMap[item.Get()] = true
	}

	// Check for common elements in the current Uint16s slice
	for _, item := range s {
		val := item.Get()
		if valueMap[val] {
			result = append(result, item)
		}
	}

	return result
}

// Difference computes the difference of the current Uint16s slice with another Uint16s slice.
// The result contains only the elements that are in the current slice but not in the other
// slice.
func (s Uint16s) Difference(other Uint16s) Uint16s {
	valueMap := make(map[uint16]bool)
	var result Uint16s

	// Mark all elements in the other Uint16s slice
	for _, item := range other {
		valueMap[item.Get()] = true
	}

	// Add elements from the current Uint16s slice that are not in the other Uint16s slice
	for _, item := range s {
		val := item.Get()
		if !valueMap[val] {
			result = append(result, item)
		}
	}

	return result
}
