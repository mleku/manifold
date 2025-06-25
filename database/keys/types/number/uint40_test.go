package number

import (
	"bytes"
	"reflect"
	"testing"
)

func TestUint40(t *testing.T) {
	// Test cases for Get
	tests := []struct {
		name        string
		value       uint64
		expectedErr bool
	}{
		{"Minimum Value", 0, false},
		{"Maximum Value", MaxUint40, false},
		{"Value in Range", 109951162777, false}, // Example value within the range
		{"Value Exceeds Range", MaxUint40 + 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codec := new(Uint40)

			// Test Set
			err := codec.Set(tt.value)
			if tt.expectedErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Test Get getter
			if codec.Get() != tt.value {
				t.Errorf("Uint40 mismatch: got %d, expected %d", codec.Get(), tt.value)
			}

			// Test MarshalWrite and UnmarshalRead
			buf := new(bytes.Buffer)

			// Marshal to a buffer
			if err = codec.MarshalWrite(buf); err != nil {
				t.Fatalf("MarshalWrite failed: %v", err)
			}

			// Validate encoded size is 5 bytes
			encoded := buf.Bytes()
			if len(encoded) != 5 {
				t.Fatalf("encoded size mismatch: got %d bytes, expected 5 bytes", len(encoded))
			}

			// Decode from the buffer
			decoded := new(Uint40)
			if err = decoded.UnmarshalRead(buf); err != nil {
				t.Fatalf("UnmarshalRead failed: %v", err)
			}

			// Validate decoded value
			if decoded.Get() != tt.value {
				t.Errorf("Decoded value mismatch: got %d, expected %d", decoded.Get(), tt.value)
			}
		})
	}
}

func TestUint40sSetOperations(t *testing.T) {
	// Helper function to create a Uint64 with a specific value
	createUint64 := func(value uint64) *Uint40 {
		u := &Uint40{}
		u.Set(value)
		return u
	}

	// Prepare test data
	a := createUint64(1)
	b := createUint64(2)
	c := createUint64(3)
	d := createUint64(4)
	e := createUint64(1) // Duplicate of a

	// Define slices
	set1 := Uint40s{a, b, c}              // [1, 2, 3]
	set2 := Uint40s{d, e, b}              // [4, 1, 2]
	expectedUnion := Uint40s{a, b, c, d}  // [1, 2, 3, 4]
	expectedIntersection := Uint40s{a, b} // [1, 2]
	expectedDifference := Uint40s{c}      // [3]

	// Test Union
	t.Run("Union", func(t *testing.T) {
		result := set1.Union(set2)
		if !reflect.DeepEqual(getUint40Values(result), getUint40Values(expectedUnion)) {
			t.Errorf("Union failed: expected %v, got %v", getUint40Values(expectedUnion), getUint40Values(result))
		}
	})

	// Test Intersection
	t.Run("Intersection", func(t *testing.T) {
		result := set1.Intersection(set2)
		if !reflect.DeepEqual(getUint40Values(result), getUint40Values(expectedIntersection)) {
			t.Errorf("Intersection failed: expected %v, got %v", getUint40Values(expectedIntersection), getUint40Values(result))
		}
	})

	// Test Difference
	t.Run("Difference", func(t *testing.T) {
		result := set1.Difference(set2)
		if !reflect.DeepEqual(getUint40Values(result), getUint40Values(expectedDifference)) {
			t.Errorf("Difference failed: expected %v, got %v", getUint40Values(expectedDifference), getUint40Values(result))
		}
	})
}

// Helper function to extract uint64 values from Uint40s
func getUint40Values(slice Uint40s) []uint64 {
	var values []uint64
	for _, item := range slice {
		values = append(values, item.Get())
	}
	return values
}
