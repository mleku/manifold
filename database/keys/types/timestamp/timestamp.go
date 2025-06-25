package timestamp

import (
	"bytes"
	"io"

	"github.com/mleku/manifold/chk"
	"github.com/mleku/manifold/database/keys/types/number"
)

const Len = 8

type T struct{ val int64 }

func (ts *T) FromInt(t int)     { ts.val = int64(t) }
func (ts *T) FromInt64(t int64) { ts.val = int64(t) }

func FromBytes(timestampBytes []byte) (ts *T, err error) {
	v := new(number.Uint64)
	if err = v.UnmarshalRead(bytes.NewBuffer(timestampBytes)); chk.E(err) {
		return
	}
	ts = &T{val: int64(v.Get())}
	return
}

func (ts *T) ToTimestamp() (timestamp int64) {
	return
}
func (ts *T) Bytes() (b []byte, err error) {
	v := new(number.Uint64)
	buf := new(bytes.Buffer)
	if err = v.MarshalWrite(buf); chk.E(err) {
		return
	}
	b = buf.Bytes()
	return
}

func (ts *T) MarshalWrite(w io.Writer) (err error) {
	v := new(number.Uint64)
	if err = v.MarshalWrite(w); chk.E(err) {
		return
	}
	return
}

func (ts *T) UnmarshalRead(r io.Reader) (err error) {
	v := new(number.Uint64)
	if err = v.UnmarshalRead(r); chk.E(err) {
		return
	}
	ts.val = int64(v.Get())
	return
}
