package indexes

import (
	"io"
	"reflect"

	"github.com/mleku/manifold/chk"
	"github.com/mleku/manifold/codec"
	"github.com/mleku/manifold/database/indexes/types/fullid"
	"github.com/mleku/manifold/database/indexes/types/fulltext"
	"github.com/mleku/manifold/database/indexes/types/identhash"
	"github.com/mleku/manifold/database/indexes/types/idhash"
	. "github.com/mleku/manifold/database/indexes/types/number"
	"github.com/mleku/manifold/database/indexes/types/pubhash"
)

type P struct {
	val []byte
}

func NewPrefix(prf ...int) (p *P) {
	if len(prf) > 0 {
		return &P{[]byte(Prefix(prf[0]))}
	} else {
		return &P{[]byte{0, 0}}
	}
}

func (p *P) Bytes() (b []byte) { return p.val }

func (p *P) MarshalWrite(w io.Writer) (err error) {
	_, err = w.Write(p.val)
	return
}

func (p *P) UnmarshalRead(r io.Reader) (err error) {
	_, err = r.Read(p.val)
	return
}

type I string

func (i I) Write(w io.Writer) (n int, err error) { return w.Write([]byte(i)) }

// Prefix returns the two byte human-readable prefixes that go in front of
// database indexes.
func Prefix(prf int) (i I) {
	switch prf {
	case Event:
		return "ev"
	case Id:
		return "id"
	case IdPubkeyTimestamp:
		return "fi"
	case PubkeyTimestamp:
		return "pt"
	case Timestamp:
		return "ts"
	case PubkeyTagTimestamp:
		return "tp"
	case TagTimestamp:
		return "tt"
	case FulltextWord:
		return "fw"
	}
	return
}

type Encs []codec.I

// T is a wrapper around an array of codec.I. The caller provides the Encs so
// they can then call the accessor function of the codec.I implementation.
type T struct {
	Encs
}

// New creates a new indexes. The helper functions below have an encode and
// decode variant, the decode variant does not add the prefix encoder because it
// has been read by Identify.
func New(encoders ...codec.I) (i *T) { return &T{encoders} }

func (t *T) MarshalWrite(w io.Writer) (err error) {
	for _, e := range t.Encs {
		if e == nil || reflect.ValueOf(e).IsNil() {
			// allow a field to be empty, as is needed for search indexes to
			// create search
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

// by eliminating kinds and all the categories of nonsense associated with them,
// the specification becomes a lot simpler. There is no "kind" in manifold; such
// data would be a tag, like mimetype, and/or encoding.

// Event is the whole event stored in binary format
//
//	[ prefix ][ 8 byte serial ] [ event in binary format ]

const Event = 0

func EventVars() (ser *Uint40) {
	ser = new(Uint40)
	return
}
func EventEnc(ser *Uint40) (enc *T) {
	return New(NewPrefix(Event), ser)
}
func EventDec(ser *Uint40) (enc *T) {
	return New(NewPrefix(), ser)
}

// Id contains a truncated 8 byte hash of an event index. This is the
// secondary key of an event, the primary key is the serial found in the
// Event.
//
// [ prefix ][ 8 bytes truncated hash of Id ][ 8 serial ]
const Id = 1

func IdVars() (id *idhash.T, ser *Uint40) {
	id = idhash.New()
	ser = new(Uint40)
	return
}
func IdEnc(id *idhash.T, ser *Uint40) (enc *T) {
	return New(NewPrefix(Id), id, ser)
}
func IdSearch(id *idhash.T) (enc *T) {
	return New(NewPrefix(Id), id)
}
func IdDec(id *idhash.T, ser *Uint40) (enc *T) {
	return New(NewPrefix(), id, ser)
}

// IdPubkeyTimestamp is an index designed to enable sorting and filtering of
// results found via other indexes, without having to decode the event.
//
// [ prefix ][ 8 bytes serial ][ 32 bytes full event ID ][ 8 bytes truncated hash of pubkey ][ 8 bytes timestamp ]
const IdPubkeyTimestamp = 2

func IdPubkeyTimestampVars() (ser *Uint40, t *fullid.T, p *pubhash.T, ts *Uint64) {
	ser = new(Uint40)
	t = fullid.New()
	p = pubhash.New()
	ts = new(Uint64)
	return
}
func IdPubkeyTimestampEnc(ser *Uint40, t *fullid.T, p *pubhash.T, ts *Uint64) (enc *T) {
	return New(NewPrefix(IdPubkeyTimestamp), ser, t, p, ts)
}
func IdPubkeyTimestampDec(ser *Uint40, t *fullid.T, p *pubhash.T, ts *Uint64) (enc *T) {
	return New(NewPrefix(), ser, t, p, ts)
}

// Timestamp is an index that allows search for the timestamp on the event.
//
// [ prefix ][ timestamp 8 bytes timestamp ][ 8 serial ]
const Timestamp = 3

func TimestampVars() (ts *Uint64, ser *Uint40) {
	ts = new(Uint64)
	ser = new(Uint40)
	return
}
func TimestampEnc(ts *Uint64, ser *Uint40) (enc *T) {
	return New(NewPrefix(Timestamp), ts, ser)
}
func TimestampDec(ts *Uint64, ser *Uint40) (enc *T) {
	return New(NewPrefix(), ts, ser)
}

// PubkeyTimestamp is a composite index that allows search by pubkey
// filtered by timestamp.
//
// [ prefix ][ 8 bytes truncated hash of pubkey ][ 8 bytes timestamp ][ 8 serial ]
const PubkeyTimestamp = 4

func PubkeyTimestampVars() (p *pubhash.T, ts *Uint64, ser *Uint40) {
	p = pubhash.New()
	ts = new(Uint64)
	ser = new(Uint40)
	return
}
func PubkeyTimestampEnc(p *pubhash.T, ts *Uint64, ser *Uint40) (enc *T) {
	return New(NewPrefix(PubkeyTimestamp), p, ts, ser)
}
func PubkeyTimestampDec(p *pubhash.T, ts *Uint64, ser *Uint40) (enc *T) {
	return New(NewPrefix(), p, ts, ser)
}

// PubkeyTagTimestamp allows searching for a pubkey, tag and timestamp.
//
// [ prefix ][ 8 bytes truncated hash of pubkey ][ 8 bytes truncated hash of key ][ 8 bytes truncated hash of value ][ 8 bytes timestamp ][ 8 serial ]
const PubkeyTagTimestamp = 5

func PubkeyTagTimestampVars() (p *pubhash.T, k, v *identhash.T, ser *Uint40) {
	k = identhash.New()
	ser = new(Uint40)
	return
}
func PubkeyTagTimestampEnc(p *pubhash.T, k, v *identhash.T, ser *Uint40) (enc *T) {
	return New(NewPrefix(PubkeyTagTimestamp), p, k, v, ser)
}
func PubkeyTagTimestampDec(p *pubhash.T, k, v *identhash.T, ser *Uint40) (enc *T) {
	return New(NewPrefix(), p, k, v, ser)
}

// TagTimestamp allows searching for a tag and filter by timestamp.
//
// [ prefix ][ 8 bytes truncated hash of key ][ 8 bytes truncated hash of value ][ 8 bytes timestamp ][ 8 serial ]
const TagTimestamp = 6

func TagTimestampVars() (k, v *identhash.T, ts *Uint64, ser *Uint40) {
	k = identhash.New()
	v = identhash.New()
	ts = new(Uint64)
	ser = new(Uint40)
	return
}
func TagTimestampEnc(k, v *identhash.T, ts *Uint64, ser *Uint40) (enc *T) {
	return New(NewPrefix(TagTimestamp), k, v, ts, ser)
}
func TagTimestampDec(k, v *identhash.T, ts *Uint64, ser *Uint40) (enc *T) {
	return New(NewPrefix(), k, v, ts, ser)
}

// FulltextWord is a fulltext word index, the index contains the whole word.
//
// [ prefix ][ full word, zero terminated ][ 3 bytes word position in content field ][ 8 serial ]
const FulltextWord = 7

func FullTextWordVars() (fw *fulltext.T, pos *Uint24, ser *Uint40) {
	fw = fulltext.New()
	pos = new(Uint24)
	ser = new(Uint40)
	return
}
func FullTextWordEnc(fw *fulltext.T, pos *Uint24, ser *Uint40) (enc *T) {
	return New(NewPrefix(FulltextWord), fw, pos, ser)
}
func FullTextWordDec(fw *fulltext.T, pos *Uint24, ser *Uint40) (enc *T) {
	return New(NewPrefix(), fw, pos, ser)
}
