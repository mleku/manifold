package fulltext

import (
	"bytes"
	"io"
)

type T struct {
	val []byte // Contains only the raw word (without the zero-byte marker)
}

func New() *T {
	return &T{}
}

// FromWord stores the word without any modifications
func (ft *T) FromWord(word []byte) {
	ft.val = word // Only store the raw word
}

// Bytes returns the raw word without any end-of-word marker
func (ft *T) Bytes() []byte {
	return ft.val
}

// MarshalWrite writes the word to the writer, appending the zero-byte marker
func (ft *T) MarshalWrite(w io.Writer) error {
	// Create a temporary buffer that contains the word and the zero-byte marker
	temp := append(ft.val, 0x00)
	_, err := w.Write(temp) // Write the buffer to the writer
	return err
}

// UnmarshalRead reads the word from the reader, stopping at the zero-byte marker
func (ft *T) UnmarshalRead(r io.Reader) error {
	var buf bytes.Buffer
	tmp := make([]byte, 1)

	// Read bytes until the zero byte is encountered
	for {
		n, err := r.Read(tmp)
		if n > 0 {
			if tmp[0] == 0x00 { // Stop on encountering the zero-byte marker
				break
			}
			buf.WriteByte(tmp[0])
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err // Handle unexpected errors
		}
	}

	// Store the raw word without the zero byte
	ft.val = buf.Bytes()
	return nil
}
