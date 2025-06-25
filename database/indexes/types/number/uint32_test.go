package number

import (
	"bytes"
	"math"
	"reflect"
	"testing"

	"lukechampine.com/frand"
)

func TestUint32(t *testing.T) {
	// Helper function to generate random 32-bit integers
	generateRandomUint32 := func() uint32 {
		return uint32(frand.Intn(math.MaxUint32)) // math.MaxUint32 == 4294967295
	}

	for i := 0; i < 100; i++ { // Run test 100 times for random values
		// Generate a random value
		randomUint32 := generateRandomUint32()
		randomInt := int(randomUint32)

		// Create a new codec
		codec := new(Uint32)

		// Test Uint32 setter and getter
		codec.Set(randomUint32)
		if codec.Get() != randomUint32 {
			t.Fatalf("Uint32 mismatch: got %d, expected %d", codec.Get(), randomUint32)
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
		if err != nil {
			t.Fatalf("MarshalWrite failed: %v", err)
		}
		encoded := bufEnc.Bytes()

		// Create a copy of encoded bytes before decoding
		bufDec := bytes.NewBuffer(encoded)

		// Decode back the value
		decoded := new(Uint32)
		err = decoded.UnmarshalRead(bufDec)
		if err != nil {
			t.Fatalf("UnmarshalRead failed: %v", err)
		}

		if decoded.Get() != randomUint32 {
			t.Fatalf("Decoded value mismatch: got %d, expected %d", decoded.Get(), randomUint32)
		}

		// Compare encoded bytes to ensure correctness
		if !bytes.Equal(encoded, bufEnc.Bytes()) {
			t.Fatalf("Byte encoding mismatch: got %v, expected %v", bufEnc.Bytes(), encoded)
		}
	}
}

func TestUint32sSetOperations(t *testing.T) {
	// Helper function to create a Uint32 with a specific value
	createUint32 := func(value uint32) *Uint32 {
		u := &Uint32{}
		u.Set(value)
		return u
	}

	// Prepare test data
	a := createUint32(1)
	b := createUint32(2)
	c := createUint32(3)
	d := createUint32(4)
	e := createUint32(1) // Duplicate of a

	// Define slices
	set1 := Uint32s{a, b, c}              // [1, 2, 3]
	set2 := Uint32s{d, e, b}              // [4, 1, 2]
	expectedUnion := Uint32s{a, b, c, d}  // [1, 2, 3, 4]
	expectedIntersection := Uint32s{a, b} // [1, 2]
	expectedDifference := Uint32s{c}      // [3]

	// Test Union
	t.Run("Union", func(t *testing.T) {
		result := set1.Union(set2)
		if !reflect.DeepEqual(getUint32Values(result), getUint32Values(expectedUnion)) {
			t.Errorf("Union failed: expected %v, got %v", getUint32Values(expectedUnion), getUint32Values(result))
		}
	})

	// Test Intersection
	t.Run("Intersection", func(t *testing.T) {
		result := set1.Intersection(set2)
		if !reflect.DeepEqual(getUint32Values(result), getUint32Values(expectedIntersection)) {
			t.Errorf("Intersection failed: expected %v, got %v", getUint32Values(expectedIntersection), getUint32Values(result))
		}
	})

	// Test Difference
	t.Run("Difference", func(t *testing.T) {
		result := set1.Difference(set2)
		if !reflect.DeepEqual(getUint32Values(result), getUint32Values(expectedDifference)) {
			t.Errorf("Difference failed: expected %v, got %v", getUint32Values(expectedDifference), getUint32Values(result))
		}
	})
}

// Helper function to extract uint64 values from Uint32s
func getUint32Values(slice Uint32s) []uint32 {
	var values []uint32
	for _, item := range slice {
		values = append(values, item.Get())
	}
	return values
}
