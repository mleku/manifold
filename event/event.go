package event

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/mleku/manifold/chk"
	"github.com/mleku/manifold/ec/schnorr"
	"github.com/mleku/manifold/errorf"
	"github.com/mleku/manifold/ints"
	"github.com/mleku/manifold/log"
)

func WriteTo(w io.Writer, counter *int, b ...[]byte) (err error) {
	if counter == nil {
		err = errorf.E("counter is nil")
		return
	}
	var n int
	for _, f := range b {
		if n, err = w.Write(f); chk.E(err) {
			return
		}
	}
	*counter += n
	return
}

type Tag struct {
	key   []byte
	value []byte
}

func (t Tag) String() string {
	return fmt.Sprintf("%s:%s", t.key, t.value)
}

type E struct {
	Pubkey    []byte
	Timestamp int64
	Content   []byte
	Tags      []Tag
	Signature []byte
}

func (e E) String() string {
	return fmt.Sprintf(
		"Pubkey (Base64): %s\nTimestamp (Decimal): %d\nContent: %s\nTags: %v\nSignature (Base64): %s",
		base64.RawStdEncoding.EncodeToString(e.Pubkey), // Base64-encoded Pubkey
		e.Timestamp,       // Decimal Timestamp
		string(e.Content), // Convert Content to string
		e.Tags,            // Use default formatting for Tags
		base64.RawStdEncoding.EncodeToString(e.Signature), // Base64-encoded Signature
	)
}

var Sentinels = [][]byte{
	[]byte("PUBKEY:"),
	[]byte("TIMESTAMP:"),
	[]byte("CONTENT:"),
	[]byte("TAG:"),
	[]byte("SIGNATURE:"),
}

func Split(b []byte) (e E, err error) {
	if !bytes.HasPrefix(b, Sentinels[0]) {
		err = errorf.E("pubkey not found")
	}
	b = b[len(Sentinels[0]):]
	// expect 64 hex bytes
	if len(b) < 44 {
		err = errorf.E("pubkey not found")
		return
	}
	pk := b[:43]
	e.Pubkey = make([]byte, schnorr.PubKeyBytesLen)
	if _, err = base64.RawStdEncoding.Decode(e.Pubkey, pk); chk.E(err) {
		err = errorf.E("pubkey did not decode correctly")
		return
	}
	b = b[43:]
	if b[0] != '\n' {
		err = errorf.E("pubkey did not decode correctly")
		return
	}
	b = b[1:]
	if !bytes.HasPrefix(b, Sentinels[1]) {
		err = errorf.E("timestamp not found")
		return
	}
	b = b[len(Sentinels[1]):]
	ts := ints.New(0)
	if b, err = ts.Unmarshal(b); chk.E(err) {
		err = errorf.E("timestamp did not decode correctly: %v", err)
		return
	}
	e.Timestamp = ts.Int64()
	if b[0] != '\n' {
		err = errorf.E("timestamp did not decode correctly")
		return
	}
	b = b[1:]
	if !bytes.HasPrefix(b, Sentinels[2]) {
		err = errorf.E("content not found")
		return
	}
	b = b[len(Sentinels[2]):]
	// after the content should be tags
	next := bytes.Index(b, Sentinels[3])
	// or if not, a signature
	if next == -1 {
		next = bytes.Index(b, Sentinels[4])
		e.Content = b[:next-1]
		b = b[next:]
		log.I.S(b)
		// next is just signature
		log.I.F("skipping to signature")
		goto sig
	}
	e.Content = b[:next]
	b = b[next:]
	for {
		b = b[len(Sentinels[3]):]
		next = bytes.Index(b, Sentinels[3])
		if next == -1 {
			next = bytes.Index(b, Sentinels[4])
			tag := b[:next-1]
			split := bytes.SplitN(tag, []byte(":"), 2)
			t := Tag{split[0], split[1]}
			e.Tags = append(e.Tags, t)
			b = b[next:]
			goto sig
		}
		tag := b[:next-1]
		split := bytes.SplitN(tag, []byte(":"), 2)
		t := Tag{split[0], split[1]}
		e.Tags = append(e.Tags, t)
		b = b[next:]
	}
sig:
	if !bytes.HasPrefix(b, Sentinels[4]) {
		err = errorf.E("signature not found")
		return
	}
	b = b[len(Sentinels[4]):]
	if len(b) < 86 {
		err = errorf.E("signature not found")
		return
	}
	sig := b[:86]
	e.Signature = make([]byte, schnorr.SignatureSize)
	log.I.F("signature: %s", sig)
	base64.RawStdEncoding.Decode(e.Signature, sig)
	// anything else is garbage
	return
}
