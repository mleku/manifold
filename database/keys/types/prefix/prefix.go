package prefix

import (
	"io"

	"github.com/mleku/manifold/database/prefixes"
)

const Len = 2

type T struct {
	val []byte
}

func New(prf ...int) (p *T) {
	if len(prf) > 0 {
		return &T{[]byte(prefixes.Prefix(prf[0]))}
	} else {
		return &T{[]byte{0, 0}}
	}
}

func (p *T) Bytes() (b []byte) { return p.val }

func (p *T) MarshalWrite(w io.Writer) (err error) {
	_, err = w.Write(p.val)
	return
}

func (p *T) UnmarshalRead(r io.Reader) (err error) {
	_, err = r.Read(p.val)
	return
}
