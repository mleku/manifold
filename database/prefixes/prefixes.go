package prefixes

import (
	"io"
)

const Len = 2

type I string

func (i I) Write(w io.Writer) (n int, err error) { return w.Write([]byte(i)) }

// Prefix returns the two byte human-readable prefixes that go in front of
// database keys.
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
	case TagTimestamp:
		return "tt"
	case PubkeyTagTimestamp:
		return "tp"
	case FulltextWord:
		return "fw"
	}
	return
}

// the prefixes.
//
// by eliminating kinds and all the categories of nonsense associated with them,
// the specification becomes a lot simpler. There is no "kind" in manifold; such
// data would be a tag, like mimetype, and/or encoding.

const (
	// Event is the whole event stored in binary format
	//
	//   [ prefix ][ 8 byte serial ] [ event in binary format ]
	Event = iota

	// Id contains a truncated 8 byte hash of an event index. This is the
	// secondary key of an event, the primary key is the serial found in the
	// Event.
	//
	// [ prefix ][ 8 bytes truncated hash of Id ][ 8 serial ]
	Id

	// IdPubkeyTimestamp is an index designed to enable sorting and filtering of
	// results found via other indexes, without having to decode the event.
	//
	// [ prefix ][ 8 bytes serial ][ 32 bytes full event ID ][ 8 bytes truncated hash of pubkey ][ 8 bytes timestamp ]
	IdPubkeyTimestamp

	// Timestamp is an index that allows search for the timestamp on the event.
	//
	// [ prefix ][ timestamp 8 bytes timestamp ][ 8 serial ]
	Timestamp

	// PubkeyTimestamp is a composite index that allows search by pubkey
	// filtered by timestamp.
	//
	// [ prefix ][ 8 bytes truncated hash of pubkey ][ 8 bytes timestamp ][ 8 serial ]
	PubkeyTimestamp

	// PubkeyTagTimestamp allows searching for a pubkey, tag and timestamp.
	//
	// [ prefix ][ 8 bytes truncated hash of pubkey ][ 8 bytes truncated hash of key ][ 8 bytes truncated hash of value ][ 8 bytes timestamp ][ 8 serial ]
	PubkeyTagTimestamp

	// TagTimestamp allows searching for a tag and filter by timestamp.
	//
	// [ prefix ][ 8 bytes truncated hash of key ][ 8 bytes truncated hash of value ][ 8 bytes timestamp ][ 8 serial ]
	TagTimestamp

	// FulltextWord is a fulltext word index, the index contains the whole word.
	//
	// [ prefix ][ full word, zero terminated ][ 3 bytes word position in content field ][ 8 serial ]
	FulltextWord
)
