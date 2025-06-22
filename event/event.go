package event

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/mleku/manifold/chk"
	"github.com/mleku/manifold/ec/schnorr"
	"github.com/mleku/manifold/errorf"
	"github.com/mleku/manifold/ints"
	"github.com/mleku/manifold/p256k"
	"github.com/mleku/manifold/sha256"
	"github.com/mleku/manifold/signer"
)

func WriteText(w io.Writer, b []byte) (err error) {
out:
	for i := range b {
		switch b[i] {
		case '\n':
			if _, err = w.Write([]byte{'\\', 'n'}); chk.E(err) {
				break out
			}
		case '\\':
			if _, err = w.Write([]byte{'\\', '\\'}); chk.E(err) {
				break out
			}
		default:
			if _, err = w.Write(b[i : i+1]); chk.E(err) {
				break out
			}
		}
	}
	return
}

func ReadText(r io.Reader) (b []byte, err error) {
	buf := new(bytes.Buffer)
	var inEscape bool
	var n int
	rb := make([]byte, 1)
	for {
		if n, err = r.Read(rb); err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		if n == 0 {
			break
		}
		if rb[0] == '\\' && !inEscape {
			inEscape = true
			// next byte determines which
			continue
		}
		if inEscape {
			switch rb[0] {
			case 'n':
				rb[0] = '\n'
			case '\\':
				rb[0] = '\\'
			}
			inEscape = false
		}
		buf.WriteByte(rb[0])
	}
	b = buf.Bytes()
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

func (e *E) String() string {
	return fmt.Sprintf(
		"Pubkey (Base64): %s\nTimestamp (Decimal): %d\nContent: %s\nTags: %v\nSignature (Base64): %s",
		base64.RawStdEncoding.EncodeToString(e.Pubkey),
		e.Timestamp,
		e.Content,
		e.Tags,
		base64.RawStdEncoding.EncodeToString(e.Signature),
	)
}

const (
	PUBKEY int = iota
	TIMESTAMP
	CONTENT
	TAG
	SIGNATURE
)

var Sentinels = [][]byte{
	[]byte("PUBKEY:"),
	[]byte("TIMESTAMP:"),
	[]byte("CONTENT:"),
	[]byte("TAG:"),
	[]byte("SIGNATURE:"),
}

func (e *E) Unmarshal(data []byte) (err error) {
	founds := make([]bool, len(Sentinels))
	scanner := bufio.NewScanner(bytes.NewBuffer(data))
	buf := make([]byte, 1_000_000)
	scanner.Buffer(buf, len(buf))
	var lines int
	for scanner.Scan() {
		if scanner.Err() != nil {
			err = scanner.Err()
			return
		}
		line := scanner.Bytes()
		lines++
		switch {
		case bytes.HasPrefix(line, Sentinels[PUBKEY]):
			if founds[PUBKEY] {
				err = errorf.E("multiple PUBKEY found at line %d\n%s", lines, data)
				return
			}
			founds[PUBKEY] = true
			e.Pubkey = make([]byte, schnorr.PubKeyBytesLen)
			if _, err = base64.RawStdEncoding.Decode(e.Pubkey, line[len(Sentinels[PUBKEY]):]); chk.E(err) {
				return
			}
		case bytes.HasPrefix(line, Sentinels[TIMESTAMP]):
			switch {
			case !founds[PUBKEY]:
				err = errorf.E("TIMESTAMP found before PUBKEY at line %d\n%s", lines, data)
				return
			case founds[TIMESTAMP]:
				err = errorf.E("multiple TIMESTAMP found at line %d\n%s", lines, data)
				return
			}
			founds[TIMESTAMP] = true
			ts := ints.New(int64(0))
			if _, err = ts.Unmarshal(line[len(Sentinels[TIMESTAMP]):]); chk.E(err) {
				return
			}
			e.Timestamp = ts.Int64()
		case bytes.HasPrefix(line, Sentinels[CONTENT]):
			switch {
			case !founds[PUBKEY]:
				err = errorf.E("CONTENT found before PUBKEY at line %d\n%s", lines, data)
				return
			case !founds[TIMESTAMP]:
				err = errorf.E("CONTENT found before TIMESTAMP at line %d\n%s", lines, data)
				return
			case founds[CONTENT]:
				err = errorf.E("multiple CONTENT found at line %d\n%s", lines, data)
				return
			}
			founds[CONTENT] = true
			content := line[len(Sentinels[CONTENT]):]
			if e.Content, err = ReadText(bytes.NewBuffer(content)); chk.E(err) {
				return
			}
		case bytes.HasPrefix(line, Sentinels[TAG]):
			switch {
			case !founds[PUBKEY]:
				err = errorf.E("TAG found before PUBKEY at line %d\n%s", lines, data)
				return
			case !founds[TIMESTAMP]:
				err = errorf.E("TAG found before TIMESTAMP at line %d\n%s", lines, data)
				return
			case !founds[CONTENT]:
				err = errorf.E("TAG found before CONTENT at line %d\n%s", lines, data)
				return
			}
			line = line[len(Sentinels[TAG]):]
			keyEnd := bytes.IndexByte(line, ':')
			if keyEnd == -1 {
				err = errorf.E("invalid TAG format\n%s", lines, data)
				return
			}
			var key []byte
			if key, err = ReadText(bytes.NewBuffer(line[:keyEnd])); chk.E(err) {
				return
			}
			var value []byte
			if value, err = ReadText(bytes.NewBuffer(line[keyEnd+1:])); chk.E(err) {
				return
			}
			e.Tags = append(e.Tags, Tag{key, value})
		case bytes.HasPrefix(line, Sentinels[SIGNATURE]):
			switch {
			case !founds[PUBKEY]:
				err = errorf.E("SIGNATURE found before PUBKEY at line %d\n%s", lines, data)
				return
			case !founds[TIMESTAMP]:
				err = errorf.E("SIGNATURE found before TIMESTAMP at line %d\n%s", lines, data)
				return
			case !founds[CONTENT]:
				err = errorf.E("SIGNATURE found before CONTENT at line %d\n%s", lines, data)
				return
			case founds[SIGNATURE]:
				err = errorf.E("multiple SIGNATURE found\n%s", lines, data)
				return
			}
			founds[SIGNATURE] = true
			e.Signature = make([]byte, schnorr.SignatureSize)
			if _, err = base64.RawStdEncoding.Decode(e.Signature, line[len(Sentinels[SIGNATURE]):]); chk.E(err) {
				return
			}
		default:
			err = errorf.E("unknown sentinel on line %d: '%s'\n%s", lines, line, data)
			return
		}
	}
	return
}

func (e *E) Marshal() (data []byte, err error) {
	buf := new(bytes.Buffer)
out:
	for i := range Sentinels {
		if i == SIGNATURE && e.Signature == nil {
			// if no signature is present, this means it should be marshaled in
			// the canonical format to be hashed to generate the message hash to
			// sign.
		} else {
			if i > 0 {
				buf.WriteByte('\n')
			}
			buf.Write(Sentinels[i])
		}
		switch i {
		case PUBKEY:
			b := make([]byte, 43)
			base64.RawStdEncoding.Encode(b, e.Pubkey)
			buf.Write(b)
		case TIMESTAMP:
			ts := ints.New(e.Timestamp)
			b := ts.Marshal(nil)
			buf.Write(b)
		case CONTENT:
			if err = WriteText(buf, e.Content); chk.E(err) {
				return
			}
		case TAG:
			for t, v := range e.Tags {
				if err = WriteText(buf, v.key); chk.E(err) {
					return
				}
				buf.WriteByte(':')
				if err = WriteText(buf, v.value); chk.E(err) {
					return
				}
				if t < len(e.Tags)-1 {
					buf.WriteByte('\n')
					buf.Write(Sentinels[i])
				}
			}
		case SIGNATURE:
			// if no signature is present, this means it should be marshaled in
			// the canonical format to be hashed to generate the message hash to
			// sign.
			if e.Signature == nil {
				break out
			}
			b := make([]byte, 86)
			base64.RawStdEncoding.Encode(b, e.Signature)
			buf.Write(b)
		}
	}
	data = buf.Bytes()
	return
}

func (e *E) Id() (id []byte, err error) {
	e2 := &E{
		Pubkey:    e.Pubkey,
		Timestamp: e.Timestamp,
		Content:   e.Content,
		Tags:      e.Tags,
	}
	var data []byte
	if data, err = e2.Marshal(); err != nil {
		return
	}
	id = sha256.Sum256Bytes(data)
	return
}

func (e *E) Sign(sign signer.I) (err error) {
	if e.Signature != nil {
		err = errorf.E("event already signed")
	}
	var id []byte
	if id, err = e.Id(); chk.E(err) {
		return
	}
	if e.Signature, err = sign.Sign(id); chk.E(err) {
		return
	}
	return
}

func (e *E) Verify() (valid bool, err error) {
	pub := new(p256k.Signer)
	if err = pub.InitPub(e.Pubkey); chk.E(err) {
		return
	}
	var id []byte
	if id, err = e.Id(); chk.E(err) {
		return
	}
	if valid, err = pub.Verify(id, e.Signature); chk.E(err) {
		return
	}
	return
}
