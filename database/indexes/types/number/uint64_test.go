package number

import (
	"bytes"
	"math"
	"reflect"
	"testing"

	"lukechampine.com/frand"
	"manifold.mleku.dev/chk"
)

func TestUint64(t *testing.T) {
	// Helper function to generate random 64-bit integers
	generateRandomUint64 := func() uint64 {
		return frand.Uint64n(math.MaxUint64) // math.MaxUint64 == 18446744073709551615
	}

	for i := 0; i < 100; i++ { // Run test 100 times for random values
		// Generate a random value
		randomUint64 := generateRandomUint64()
		randomInt := int(randomUint64)

		// Create a new codec
		codec := new(Uint64)

		// Test UInt64 setter and getter
		codec.Set(randomUint64)
		if codec.Get() != randomUint64 {
			t.Fatalf("Uint64 mismatch: got %d, expected %d", codec.Get(), randomUint64)
		}

		// Test GetInt setter and getter
		codec.SetInt(randomInt)
		if codec.Int() != randomInt {
			t.Fatalf("GetInt mismatch: got %d, expected %d", codec.Int(), randomInt)
		}

		// Test encoding to []byte and decoding back
		bufEnc := new(bytes.Buffer)

		// MarshalWrite
		err := codec.MarshalWrite(bufEnc)
		if chk.E(err) {
			t.Fatalf("MarshalWrite failed: %v", err)
		}
		encoded := bufEnc.Bytes()

		// Create a buffer for decoding
		bufDec := bytes.NewBuffer(encoded)

		// Decode back the value
		decoded := new(Uint64)
		err = decoded.UnmarshalRead(bufDec)
		if chk.E(err) {
			t.Fatalf("UnmarshalRead failed: %v", err)
		}

		if decoded.Get() != randomUint64 {
			t.Fatalf("Decoded value mismatch: got %d, expected %d", decoded.Get(), randomUint64)
		}

		// Compare encoded bytes to ensure correctness
		if !bytes.Equal(encoded, bufEnc.Bytes()) {
			t.Fatalf("Byte encoding mismatch: got %v, expected %v", bufEnc.Bytes(), encoded)
		}
	}
}

func TestUint64sSetOperations(t *testing.T) {
	// Helper function to create a Uint64 with a specific value
	createUint64 := func(value uint64) *Uint64 {
		u := &Uint64{}
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
	set1 := Uint64s{a, b, c}              // [1, 2, 3]
	set2 := Uint64s{d, e, b}              // [4, 1, 2]
	expectedUnion := Uint64s{a, b, c, d}  // [1, 2, 3, 4]
	expectedIntersection := Uint64s{a, b} // [1, 2]
	expectedDifference := Uint64s{c}      // [3]

	// Test Union
	t.Run("Union", func(t *testing.T) {
		result := set1.Union(set2)
		if !reflect.DeepEqual(getUint64Values(result), getUint64Values(expectedUnion)) {
			t.Errorf("Union failed: expected %v, got %v", getUint64Values(expectedUnion), getUint64Values(result))
		}
	})

	// Test Intersection
	t.Run("Intersection", func(t *testing.T) {
		result := set1.Intersection(set2)
		if !reflect.DeepEqual(getUint64Values(result), getUint64Values(expectedIntersection)) {
			t.Errorf("Intersection failed: expected %v, got %v", getUint64Values(expectedIntersection), getUint64Values(result))
		}
	})

	// Test Difference
	t.Run("Difference", func(t *testing.T) {
		result := set1.Difference(set2)
		if !reflect.DeepEqual(getUint64Values(result), getUint64Values(expectedDifference)) {
			t.Errorf("Difference failed: expected %v, got %v", getUint64Values(expectedDifference), getUint64Values(result))
		}
	})
}

// Helper function to extract uint64 values from Uint64s
func getUint64Values(slice Uint64s) []uint64 {
	var values []uint64
	for _, item := range slice {
		values = append(values, item.Get())
	}
	return values
}
