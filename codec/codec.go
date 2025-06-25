package codec

import (
	"io"
)

type I interface {
	MarshalWrite(w io.Writer) (err error)
	UnmarshalRead(r io.Reader) (err error)
}
