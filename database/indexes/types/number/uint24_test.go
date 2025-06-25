package number

import (
	"bytes"
	"reflect"
	"testing"

	"manifold.mleku.dev/chk"
)

func TestUint24(t *testing.T) {
	tests := []struct {
		name        string
		value       uint32
		expectedErr bool
	}{
		{"Minimum Value", 0, false},
		{"Maximum Value", MaxUint24, false},
		{"Value in Range", 8374263, false},           // Example value within the range
		{"Value Exceeds Range", MaxUint24 + 1, true}, // Exceeds 24-bit limit
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codec := new(Uint24)

			// Test Set
			err := codec.Set(tt.value)
			if tt.expectedErr {
				if !chk.E(err) {
					t.Errorf("expected error but got none")
				}
				return
			} else if chk.E(err) {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Test Get getter
			if codec.Get() != tt.value {
				t.Errorf("Get mismatch: got %d, expected %d", codec.Get(), tt.value)
			}

			// Test MarshalWrite and UnmarshalRead
			buf := new(bytes.Buffer)

			// MarshalWrite directly to the buffer
			if err := codec.MarshalWrite(buf); chk.E(err) {
				t.Fatalf("MarshalWrite failed: %v", err)
			}

			// Validate encoded size is 3 bytes
			encoded := buf.Bytes()
			if len(encoded) != 3 {
				t.Fatalf("encoded size mismatch: got %d bytes, expected 3 bytes", len(encoded))
			}

			// Decode from the buffer
			decoded := new(Uint24)
			if err := decoded.UnmarshalRead(buf); chk.E(err) {
				t.Fatalf("UnmarshalRead failed: %v", err)
			}

			// Validate decoded value
			if decoded.Get() != tt.value {
				t.Errorf("Decoded value mismatch: got %d, expected %d", decoded.Get(), tt.value)
			}
		})
	}
}

func TestUint24sSetOperations(t *testing.T) {
	// Helper function to create a Uint24 with a specific value
	createUint24 := func(value uint32) *Uint24 {
		u := &Uint24{}
		u.Set(value)
		return u
	}

	// Prepare test data
	a := createUint24(1)
	b := createUint24(2)
	c := createUint24(3)
	d := createUint24(4)
	e := createUint24(1) // Duplicate of a

	// Define slices
	set1 := Uint24s{a, b, c}              // [1, 2, 3]
	set2 := Uint24s{d, e, b}              // [4, 1, 2]
	expectedUnion := Uint24s{a, b, c, d}  // [1, 2, 3, 4]
	expectedIntersection := Uint24s{a, b} // [1, 2]
	expectedDifference := Uint24s{c}      // [3]

	// Test Union
	t.Run("Union", func(t *testing.T) {
		result := set1.Union(set2)
		if !reflect.DeepEqual(getUint24Values(result), getUint24Values(expectedUnion)) {
			t.Errorf("Union failed: expected %v, got %v", getUint24Values(expectedUnion), getUint24Values(result))
		}
	})

	// Test Intersection
	t.Run("Intersection", func(t *testing.T) {
		result := set1.Intersection(set2)
		if !reflect.DeepEqual(getUint24Values(result), getUint24Values(expectedIntersection)) {
			t.Errorf("Intersection failed: expected %v, got %v", getUint24Values(expectedIntersection), getUint24Values(result))
		}
	})

	// Test Difference
	t.Run("Difference", func(t *testing.T) {
		result := set1.Difference(set2)
		if !reflect.DeepEqual(getUint24Values(result), getUint24Values(expectedDifference)) {
			t.Errorf("Difference failed: expected %v, got %v", getUint24Values(expectedDifference), getUint24Values(result))
		}
	})
}

// Helper function to extract uint64 values from Uint24s
func getUint24Values(slice Uint24s) []uint32 {
	var values []uint32
	for _, item := range slice {
		values = append(values, item.Get())
	}
	return values
}
