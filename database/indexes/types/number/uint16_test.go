package number

import (
	"bytes"
	"math"
	"reflect"
	"testing"

	"lukechampine.com/frand"
	"manifold.mleku.dev/chk"
)

func TestUint16(t *testing.T) {
	// Helper function to generate random 16-bit integers
	generateRandomUint16 := func() uint16 {
		return uint16(frand.Intn(math.MaxUint16)) // math.MaxUint16 == 65535
	}

	for i := 0; i < 100; i++ { // Run test 100 times for random values
		// Generate a random value
		randomUint16 := generateRandomUint16()
		randomInt := int(randomUint16)

		// Create a new encodedUint16
		encodedUint16 := new(Uint16)

		// Test UInt16 setter and getter
		encodedUint16.Set(randomUint16)
		if encodedUint16.Get() != randomUint16 {
			t.Fatalf("Get mismatch: got %d, expected %d", encodedUint16.Get(), randomUint16)
		}

		// Test GetInt setter and getter
		encodedUint16.SetInt(randomInt)
		if encodedUint16.GetInt() != randomInt {
			t.Fatalf("GetInt mismatch: got %d, expected %d", encodedUint16.GetInt(), randomInt)
		}

		// Test encoding to []byte and decoding back
		bufEnc := new(bytes.Buffer)

		// MarshalWrite
		err := encodedUint16.MarshalWrite(bufEnc)
		if chk.E(err) {
			t.Fatalf("MarshalWrite failed: %v", err)
		}
		encoded := bufEnc.Bytes()

		// Create a copy of encoded bytes before decoding
		bufDec := bytes.NewBuffer(encoded)

		// Decode back the value
		decodedUint16 := new(Uint16)
		err = decodedUint16.UnmarshalRead(bufDec)
		if chk.E(err) {
			t.Fatalf("UnmarshalRead failed: %v", err)
		}

		if decodedUint16.Get() != randomUint16 {
			t.Fatalf("Decoded value mismatch: got %d, expected %d", decodedUint16.Get(), randomUint16)
		}

		// Compare encoded bytes to ensure correctness
		if !bytes.Equal(encoded, bufEnc.Bytes()) {
			t.Fatalf("Byte encoding mismatch: got %v, expected %v", bufEnc.Bytes(), encoded)
		}
	}
}

func TestUint16sSetOperations(t *testing.T) {
	// Helper function to create a Uint16 with a specific value
	createUint16 := func(value uint16) *Uint16 {
		u := &Uint16{}
		u.Set(value)
		return u
	}

	// Prepare test data
	a := createUint16(1)
	b := createUint16(2)
	c := createUint16(3)
	d := createUint16(4)
	e := createUint16(1) // Duplicate of a

	// Define slices
	set1 := Uint16s{a, b, c}              // [1, 2, 3]
	set2 := Uint16s{d, e, b}              // [4, 1, 2]
	expectedUnion := Uint16s{a, b, c, d}  // [1, 2, 3, 4]
	expectedIntersection := Uint16s{a, b} // [1, 2]
	expectedDifference := Uint16s{c}      // [3]

	// Test Union
	t.Run("Union", func(t *testing.T) {
		result := set1.Union(set2)
		if !reflect.DeepEqual(getUint16Values(result), getUint16Values(expectedUnion)) {
			t.Errorf("Union failed: expected %v, got %v", getUint16Values(expectedUnion), getUint16Values(result))
		}
	})

	// Test Intersection
	t.Run("Intersection", func(t *testing.T) {
		result := set1.Intersection(set2)
		if !reflect.DeepEqual(getUint16Values(result), getUint16Values(expectedIntersection)) {
			t.Errorf("Intersection failed: expected %v, got %v", getUint16Values(expectedIntersection), getUint16Values(result))
		}
	})

	// Test Difference
	t.Run("Difference", func(t *testing.T) {
		result := set1.Difference(set2)
		if !reflect.DeepEqual(getUint16Values(result), getUint16Values(expectedDifference)) {
			t.Errorf("Difference failed: expected %v, got %v", getUint16Values(expectedDifference), getUint16Values(result))
		}
	})
}

// Helper function to extract uint64 values from Uint16s
func getUint16Values(slice Uint16s) []uint16 {
	var values []uint16
	for _, item := range slice {
		values = append(values, item.Get())
	}
	return values
}
