package event

import (
	"io"

	"manifold.mleku.dev/chk"
	"manifold.mleku.dev/ec/schnorr"
	"manifold.mleku.dev/errorf"
	. "manifold.mleku.dev/varint"
)

var ck = chk.E
var ef = errorf.E

func (e *E) WriteBinary(w io.Writer) (err error) {
	if e == nil {
		return ef("cannot marshal nil event")
	}
	if len(e.Pubkey) != schnorr.PubKeyBytesLen {
		return ef("invalid pubkey length")
	}
	if _, err = w.Write(e.Pubkey); ck(err) {
		return
	}
	Encode(w, e.Timestamp)
	Encode(w, len(e.Content))
	if _, err = w.Write(e.Content); ck(err) {
		return
	}
	if e.Tags == nil || len(*e.Tags) == 0 {
		Encode(w, 0)
	} else {
		Encode(w, len(*e.Tags))
		for _, v := range *e.Tags {
			Encode(w, len(v.Key))
			if _, err = w.Write(v.Key); ck(err) {
				return
			}
			Encode(w, len(v.Value))
			if _, err = w.Write(v.Value); ck(err) {
				return
			}
		}
	}
	if _, err = w.Write(e.Signature); ck(err) {
		return
	}
	return
}

func (e *E) ReadBinary(r io.Reader) (err error) {
	if e == nil {
		err = ef("cannot unmarshal nil event")
		return
	}
	// read in pubkey
	e.Pubkey = make([]byte, schnorr.PubKeyBytesLen)
	if _, err = r.Read(e.Pubkey); ck(err) {
		return
	}
	var vi uint64
	// read in timestamp
	if vi, err = Decode(r); ck(err) {
		return
	}
	e.Timestamp = int64(vi)
	// read in content length
	if vi, err = Decode(r); ck(err) {
		return
	}
	// read in content
	e.Content = make([]byte, vi)
	if _, err = r.Read(e.Content); ck(err) {
		return
	}
	// read tags length
	if vi, err = Decode(r); ck(err) {
		return
	}
	for range vi {
		// read key length
		if vi, err = Decode(r); ck(err) {
			return
		}
		key := make([]byte, vi)
		// read key
		if _, err = r.Read(key); ck(err) {
			return
		}
		// read value length
		if vi, err = Decode(r); ck(err) {
			return
		}
		val := make([]byte, vi)
		// read value
		if _, err = r.Read(val); ck(err) {
			return
		}
		if e.Tags == nil {
			e.Tags = &Tags{}
		}
		*e.Tags = append(*e.Tags, Tag{key, val})
	}
	e.Signature = make([]byte, schnorr.SignatureSize)
	if _, err = r.Read(e.Signature); ck(err) {
		return
	}
	return
}
