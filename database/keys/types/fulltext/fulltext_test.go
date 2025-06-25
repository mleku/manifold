package fulltext_test

import (
	"bytes"
	"testing"

	"github.com/mleku/manifold/database/keys/types/fulltext"
)

func TestT(t *testing.T) {
	// Test cases: each contains inputs, expected serialized output, and expected result after deserialization.
	tests := []struct {
		word            []byte // Input word
		expectedBytes   []byte // Expected output from Bytes() (raw word)
		expectedEncoded []byte // Expected serialized (MarshalWrite) output (word + 0x00)
	}{
		{[]byte("example"), []byte("example"), []byte("example\x00")},
		{[]byte("golang"), []byte("golang"), []byte("golang\x00")},
		{[]byte(""), []byte(""), []byte("\x00")}, // Edge case: empty word
		{[]byte("123"), []byte("123"), []byte("123\x00")},
	}

	for _, tt := range tests {
		// Create a new object and set the word
		ft := fulltext.New()
		ft.FromWord(tt.word)

		// Ensure Bytes() returns the correct raw word
		if got := ft.Bytes(); !bytes.Equal(tt.expectedBytes, got) {
			t.Errorf("FromWord/Bytes failed: expected %q, got %q", tt.expectedBytes, got)
		}

		// Test MarshalWrite
		var buf bytes.Buffer
		if err := ft.MarshalWrite(&buf); err != nil {
			t.Fatalf("MarshalWrite failed: %v", err)
		}

		// Ensure the serialized output matches expectedEncoded
		if got := buf.Bytes(); !bytes.Equal(tt.expectedEncoded, got) {
			t.Errorf("MarshalWrite failed: expected %q, got %q", tt.expectedEncoded, got)
		}

		// Test UnmarshalRead
		newFt := fulltext.New()
		if err := newFt.UnmarshalRead(&buf); err != nil {
			t.Fatalf("UnmarshalRead failed: %v", err)
		}

		// Ensure the word after decoding matches the original word
		if got := newFt.Bytes(); !bytes.Equal(tt.expectedBytes, got) {
			t.Errorf("UnmarshalRead failed: expected %q, got %q", tt.expectedBytes, got)
		}
	}
}

func TestUnmarshalReadHandlesMissingZeroByte(t *testing.T) {
	// Special case: what happens if the zero-byte marker is missing?
	data := []byte("incomplete") // No zero-byte at the end
	reader := bytes.NewReader(data)

	ft := fulltext.New()
	err := ft.UnmarshalRead(reader)

	// Expect an EOF or similar handling
	if err == nil {
		t.Errorf("UnmarshalRead should fail gracefully on missing zero-byte, but it didn't")
	}

	// Ensure no data is stored in ft.val if no valid end-marker was encountered
	if got := ft.Bytes(); len(got) != 0 {
		t.Errorf("UnmarshalRead stored incomplete data: got %q, expected empty", got)
	}
}
