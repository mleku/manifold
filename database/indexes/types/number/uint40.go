package number

import (
	"errors"
	"io"

	"manifold.mleku.dev/chk"
)

// MaxUint40 is the maximum value of a 40-bit unsigned integer: 2^40 - 1.
const MaxUint40 uint64 = 1<<40 - 1

// Uint40 is a codec for encoding and decoding 40-bit unsigned integers.
type Uint40 struct{ value uint64 }

// Set sets the value as a 40-bit unsigned integer.
// If the value exceeds the maximum allowable value for 40 bits, it returns an error.
func (c *Uint40) Set(value uint64) error {
	if value > MaxUint40 {
		return errors.New("value exceeds 40-bit range")
	}
	c.value = value
	return nil
}

// Get gets the value as a 40-bit unsigned integer.
func (c *Uint40) Get() uint64 { return c.value }

// SetInt sets the value as an int, converting it to a 40-bit unsigned integer.
// If the value is out of the 40-bit range, it returns an error.
func (c *Uint40) SetInt(value int) error {
	if value < 0 || uint64(value) > MaxUint40 {
		return errors.New("value exceeds 40-bit range")
	}
	c.value = uint64(value)
	return nil
}

// GetInt gets the value as an int, converted from the 40-bit unsigned integer.
// Note: If the value exceeds the int range, it will be truncated.
func (c *Uint40) GetInt() int { return int(c.value) }

// MarshalWrite encodes the 40-bit unsigned integer and writes it to the provided writer.
// The encoding uses 5 bytes in BigEndian order.
func (c *Uint40) MarshalWrite(w io.Writer) (err error) {
	if c.value > MaxUint40 {
		return errors.New("value exceeds 40-bit range")
	}
	// Buffer for the 5 bytes
	buf := make([]byte, 5)
	// Write the upper 5 bytes (ignoring the most significant 3 bytes of uint64)
	buf[0] = byte((c.value >> 32) & 0xFF) // Most significant byte
	buf[1] = byte((c.value >> 24) & 0xFF)
	buf[2] = byte((c.value >> 16) & 0xFF)
	buf[3] = byte((c.value >> 8) & 0xFF)
	buf[4] = byte(c.value & 0xFF) // Least significant byte
	_, err = w.Write(buf)
	return err
}

// UnmarshalRead reads 5 bytes from the provided reader and decodes it into a 40-bit unsigned integer.
func (c *Uint40) UnmarshalRead(r io.Reader) (err error) {
	// Buffer for the 5 bytes
	buf := make([]byte, 5)
	_, err = r.Read(buf)
	if chk.E(err) {
		return err
	}
	// Decode the 5 bytes into a 40-bit unsigned integer
	c.value = (uint64(buf[0]) << 32) |
		(uint64(buf[1]) << 24) |
		(uint64(buf[2]) << 16) |
		(uint64(buf[3]) << 8) |
		uint64(buf[4])

	return nil
}

type Uint40s []*Uint40

// Union computes the union of the current Uint40s slice with another Uint40s slice. The result
// contains all unique elements from both slices.
func (s Uint40s) Union(other Uint40s) Uint40s {
	valueMap := make(map[uint64]bool)
	var result Uint40s

	// Add elements from the current Uint40s slice to the result
	for _, item := range s {
		val := item.Get()
		if !valueMap[val] {
			valueMap[val] = true
			result = append(result, item)
		}
	}

	// Add elements from the other Uint40s slice to the result
	for _, item := range other {
		val := item.Get()
		if !valueMap[val] {
			valueMap[val] = true
			result = append(result, item)
		}
	}

	return result
}

// Intersection computes the intersection of the current Uint40s slice with another Uint40s
// slice. The result contains only the elements that exist in both slices.
func (s Uint40s) Intersection(other Uint40s) Uint40s {
	valueMap := make(map[uint64]bool)
	var result Uint40s

	// Add all elements from the other Uint40s slice to the map
	for _, item := range other {
		valueMap[item.Get()] = true
	}

	// Check for common elements in the current Uint40s slice
	for _, item := range s {
		val := item.Get()
		if valueMap[val] {
			result = append(result, item)
		}
	}

	return result
}

// Difference computes the difference of the current Uint40s slice with another Uint40s slice.
// The result contains only the elements that are in the current slice but not in the other
// slice.
func (s Uint40s) Difference(other Uint40s) Uint40s {
	valueMap := make(map[uint64]bool)
	var result Uint40s

	// Mark all elements in the other Uint40s slice
	for _, item := range other {
		valueMap[item.Get()] = true
	}

	// Add elements from the current Uint40s slice that are not in the other Uint40s slice
	for _, item := range s {
		val := item.Get()
		if !valueMap[val] {
			result = append(result, item)
		}
	}

	return result
}
