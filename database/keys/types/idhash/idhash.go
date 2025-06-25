package idhash

import (
	"encoding/base64"
	"io"

	"github.com/mleku/manifold/sha256"

	"github.com/mleku/manifold/chk"
	"github.com/mleku/manifold/errorf"
	"github.com/mleku/manifold/hex"
)

const Len = 8

type T struct{ val []byte }

func New() (i *T) { return &T{make([]byte, Len)} }

func (i *T) FromId(id []byte) (err error) {
	if len(id) != sha256.Size {
		err = errorf.E("invalid Id length, got %d require %d", len(id), sha256.Size)
		return
	}
	i.val = sha256.Sum256Bytes(id)[:Len]
	return
}

func (i *T) FromIdBase64(idb64 string) (err error) {
	id := make([]byte, 43)
	if _, err = base64.RawURLEncoding.Decode(id, []byte(idb64)); chk.E(err) {
		return
	}
	i.val = sha256.Sum256Bytes(id)[:Len]
	return

}

func (i *T) FromIdHex(idh string) (err error) {
	var id []byte
	if id, err = hex.Dec(idh); chk.E(err) {
		return
	}
	if len(id) != sha256.Size {
		err = errorf.E("invalid Id length, got %d require %d", len(id), sha256.Size)
		return
	}
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
