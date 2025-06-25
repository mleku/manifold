package keys

import (
	"io"
	"reflect"

	"github.com/mleku/manifold/chk"
	"github.com/mleku/manifold/codec"
	"github.com/mleku/manifold/database/keys/types/fullid"
	"github.com/mleku/manifold/database/keys/types/fulltext"
	"github.com/mleku/manifold/database/keys/types/identhash"
	"github.com/mleku/manifold/database/keys/types/idhash"
	. "github.com/mleku/manifold/database/keys/types/number"
	"github.com/mleku/manifold/database/keys/types/prefix"
	"github.com/mleku/manifold/database/keys/types/pubhash"
	"github.com/mleku/manifold/database/prefixes"
)

type Encs []codec.I

// T is a wrapper around an array of codec.I. The caller provides the Encs so
// they can then call the accessor function of the codec.I implementation.
type T struct {
	Encs
}

// New creates a new indexes. The helper functions below have an encode and
// decode variant, the decode variant does not add the prefix encoder because it
// has been read by prefixes.Identify.
func New(encoders ...codec.I) (i *T) { return &T{encoders} }

func (t *T) MarshalWrite(w io.Writer) (err error) {
	for _, e := range t.Encs {
		if e == nil || reflect.ValueOf(e).IsNil() {
			// allow a field to be empty, as is needed for search indexes to
			// create search prefixes.
			return
		}
		if err = e.MarshalWrite(w); chk.E(err) {
			return
		}
	}
	return
}

func (t *T) UnmarshalRead(r io.Reader) (err error) {
	for _, e := range t.Encs {
		if err = e.UnmarshalRead(r); chk.E(err) {
			return
		}
	}
	return
}

func EventVars() (ser *Uint40) {
	ser = new(Uint40)
	return
}
func EventEnc(ser *Uint40) (enc *T) {
	return New(prefix.New(prefixes.Event), ser)
}
func EventDec(ser *Uint40) (enc *T) {
	return New(prefix.New(), ser)
}

func IdVars() (id *idhash.T, ser *Uint40) {
	id = idhash.New()
	ser = new(Uint40)
	return
}
func IdEnc(id *idhash.T, ser *Uint40) (enc *T) {
	return New(prefix.New(prefixes.Id), id, ser)
}
func IdSearch(id *idhash.T) (enc *T) {
	return New(prefix.New(prefixes.Id), id)
}
func IdDec(id *idhash.T, ser *Uint40) (enc *T) {
	return New(prefix.New(), id, ser)
}

type IdPubkeyTimestamp struct {
	Ser       *Uint40
	Id        *fullid.T
	Pubkey    *pubhash.T
	Kind      *Uint16
	Timestamp *Uint64
}

func IdPubkeyTimestampVars() (ser *Uint40, t *fullid.T, p *pubhash.T, ts *Uint64) {
	ser = new(Uint40)
	t = fullid.New()
	p = pubhash.New()
	ts = new(Uint64)
	return
}
func IdPubkeyTimestampEnc(ser *Uint40, t *fullid.T, p *pubhash.T, ts *Uint64) (enc *T) {
	return New(prefix.New(prefixes.IdPubkeyTimestamp), ser, t, p, ts)
}
func IdPubkeyTimestampDec(ser *Uint40, t *fullid.T, p *pubhash.T, ts *Uint64) (enc *T) {
	return New(prefix.New(), ser, t, p, ts)
}

func TimestampVars() (ts *Uint64, ser *Uint40) {
	ts = new(Uint64)
	ser = new(Uint40)
	return
}
func TimestampEnc(ts *Uint64, ser *Uint40) (enc *T) {
	return New(prefix.New(prefixes.Timestamp), ts, ser)
}
func TimestampDec(ts *Uint64, ser *Uint40) (enc *T) {
	return New(prefix.New(), ts, ser)
}

func PubkeyTimestampVars() (p *pubhash.T, ts *Uint64, ser *Uint40) {
	p = pubhash.New()
	ts = new(Uint64)
	ser = new(Uint40)
	return
}
func PubkeyTimestampEnc(p *pubhash.T, ts *Uint64, ser *Uint40) (enc *T) {
	return New(prefix.New(prefixes.PubkeyTimestamp), p, ts, ser)
}
func PubkeyTimestampDec(p *pubhash.T, ts *Uint64, ser *Uint40) (enc *T) {
	return New(prefix.New(), p, ts, ser)
}

func PubkeyTagTimestampVars() (p *pubhash.T, k, v *identhash.T, ser *Uint40) {
	k = identhash.New()
	ser = new(Uint40)
	return
}
func PubkeyTagTimestampEnc(p *pubhash.T, k, v *identhash.T, ser *Uint40) (enc *T) {
	return New(prefix.New(prefixes.PubkeyTagTimestamp), p, k, v, ser)
}
func PubkeyTagTimestampDec(p *pubhash.T, k, v *identhash.T, ser *Uint40) (enc *T) {
	return New(prefix.New(), p, k, v, ser)
}

func TagTimestampVars() (k, v *identhash.T, ts *Uint64, ser *Uint40) {
	k = identhash.New()
	v = identhash.New()
	ts = new(Uint64)
	ser = new(Uint40)
	return
}
func TagTimestampEnc(k, v *identhash.T, ts *Uint64, ser *Uint40) (enc *T) {
	return New(prefix.New(prefixes.TagTimestamp), k, v, ts, ser)
}
func TagTimestampDec(k, v *identhash.T, ts *Uint64, ser *Uint40) (enc *T) {
	return New(prefix.New(), k, v, ts, ser)
}

func FullTextWordVars() (fw *fulltext.T, pos *Uint24, ser *Uint40) {
	fw = fulltext.New()
	pos = new(Uint24)
	ser = new(Uint40)
	return
}
func FullTextWordEnc(fw *fulltext.T, pos *Uint24, ser *Uint40) (enc *T) {
	return New(prefix.New(prefixes.FulltextWord), fw, pos, ser)
}
func FullTextWordDec(fw *fulltext.T, pos *Uint24, ser *Uint40) (enc *T) {
	return New(prefix.New(), fw, pos, ser)
}
