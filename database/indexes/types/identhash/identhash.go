package identhash

import (
	"io"

	"manifold.mleku.dev/sha256"
)

const Len = 8

type T struct{ val []byte }

func New() (i *T) { return &T{make([]byte, Len)} }

func (i *T) FromIdent(id []byte) (err error) {
	i.val = sha256.Sum256Bytes(id)[:Len]
	return
}

func (i *T) Bytes() (b []byte) { return i.val }

func (i *T) MarshalWrite(w io.Writer) (err error) {
	_, err = w.Write(i.val)
	return
}

func (i *T) UnmarshalRead(r io.Reader) (err error) {
	if len(i.val) < Len {
		i.val = make([]byte, Len)
	} else {
		i.val = i.val[:Len]
	}
	_, err = r.Read(i.val)
	return
}
