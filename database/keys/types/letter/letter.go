package letter

import (
	"io"
)

const Len = 1

type T struct {
	val []byte
}

func New(letter byte) (p *T) { return &T{[]byte{letter}} }

func (p *T) Set(lb byte) { p.val = []byte{lb} }

func (p *T) Letter() byte { return p.val[0] }

func (p *T) MarshalWrite(w io.Writer) (err error) {
	_, err = w.Write(p.val)
	return
}

func (p *T) UnmarshalRead(r io.Reader) (err error) {
	val := make([]byte, 1)
	_, err = r.Read(val)
	p.val = val
	return
}
