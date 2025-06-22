package event

import (
	"bufio"
	"bytes"
	"encoding/base64"
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

type E struct {
	Pubkey    []byte
	Timestamp int64
	Content   []byte
	Tags      []Tag
	Signature []byte
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

var BinPrefix = []byte("BIN:")

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
				err = errorf.E(
					"multiple PUBKEY found at line %d\n%s",
					lines, data)
				return
			}
			founds[PUBKEY] = true
			e.Pubkey = make([]byte, schnorr.PubKeyBytesLen)
			if _, err = base64.RawURLEncoding.Decode(e.Pubkey,
				line[len(Sentinels[PUBKEY]):]); chk.E(err) {
				return
			}
		case bytes.HasPrefix(line, Sentinels[TIMESTAMP]):
			switch {
			case !founds[PUBKEY]:
				err = errorf.E(
					"TIMESTAMP found before PUBKEY at line %d\n%s",
					lines, data)
				return
			case founds[TIMESTAMP]:
				err = errorf.E(
					"multiple TIMESTAMP found at line %d\n%s",
					lines, data)
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
				err = errorf.E(
					"CONTENT found before PUBKEY at line %d\n%s",
					lines, data)
				return
			case !founds[TIMESTAMP]:
				err = errorf.E(
					"CONTENT found before TIMESTAMP at line %d\n%s",
					lines, data)
				return
			case founds[CONTENT]:
				err = errorf.E(
					"multiple CONTENT found at line %d\n%s",
					lines, data)
				return
			}
			founds[CONTENT] = true
			rawValue := line[len(Sentinels[CONTENT]):]
			var content []byte
			if bytes.HasPrefix(rawValue, []byte("b64:")) {
				// Handle Base64 decoding
				rawValue = rawValue[len("b64:"):] // Remove b64: prefix
				content = make([]byte, base64.URLEncoding.
					DecodedLen(len(rawValue))+len(BinPrefix))
				copy(content, BinPrefix)
				if _, err = base64.URLEncoding.Decode(content[len(BinPrefix):],
					rawValue); chk.E(err) {
					return
				}
			} else {
				// Handle plain text
				if content, err = ReadText(bytes.NewBuffer(rawValue)); chk.E(err) {
					return
				}
			}
			e.Content = content
		case bytes.HasPrefix(line, Sentinels[TAG]):
			switch {
			case !founds[PUBKEY]:
				err = errorf.E(
					"TAG found before PUBKEY at line %d\n%s",
					lines, data)
				return
			case !founds[TIMESTAMP]:
				err = errorf.E(
					"TAG found before TIMESTAMP at line %d\n%s",
					lines, data)
				return
			case !founds[CONTENT]:
				err = errorf.E(
					"TAG found before CONTENT at line %d\n%s",
					lines, data)
				return
			}
			line = line[len(Sentinels[TAG]):]
			keyEnd := bytes.IndexByte(line, ':')
			if keyEnd == -1 {
				err = errorf.E("invalid TAG format at line %d\n%s",
					lines, data)
				return
			}
			var key, value []byte
			if key, err = ReadText(bytes.NewBuffer(line[:keyEnd])); chk.E(err) {
				return
			}
			rawValue := line[keyEnd+1:]
			if bytes.HasPrefix(rawValue, []byte("b64:")) {
				// Handle Base64 decoding
				rawValue = rawValue[len("b64:"):] // Remove b64: prefix
				value = make([]byte, base64.URLEncoding.
					DecodedLen(len(rawValue))+len(BinPrefix))
				copy(value, BinPrefix)
				if _, err = base64.URLEncoding.Decode(value[len(BinPrefix):],
					rawValue); chk.E(err) {
					return
				}
			} else {
				// Handle plain text
				if value, err = ReadText(bytes.NewBuffer(rawValue)); chk.E(err) {
					return
				}
			}
			e.Tags = append(e.Tags, Tag{key, value})
		case bytes.HasPrefix(line, Sentinels[SIGNATURE]):
			switch {
			case !founds[PUBKEY]:
				err = errorf.E(
					"SIGNATURE found before PUBKEY at line %d\n%s",
					lines, data)
				return
			case !founds[TIMESTAMP]:
				err = errorf.E(
					"SIGNATURE found before TIMESTAMP at line %d\n%s",
					lines, data)
				return
			case !founds[CONTENT]:
				err = errorf.E(
					"SIGNATURE found before CONTENT at line %d\n%s",
					lines, data)
				return
			case founds[SIGNATURE]:
				err = errorf.E(
					"multiple SIGNATURE found\n%s",
					lines, data)
				return
			}
			founds[SIGNATURE] = true
			e.Signature = make([]byte, schnorr.SignatureSize)
			if _, err = base64.RawURLEncoding.Decode(e.Signature,
				line[len(Sentinels[SIGNATURE]):]); chk.E(err) {
				return
			}
		default:
			err = errorf.E("unknown sentinel on line %d: '%s'\n%s",
				lines, line, data)
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
			base64.RawURLEncoding.Encode(b, e.Pubkey)
			buf.Write(b)
		case TIMESTAMP:
			ts := ints.New(e.Timestamp)
			b := ts.Marshal(nil)
			buf.Write(b)
		case CONTENT:
			if bytes.HasPrefix(e.Content, BinPrefix) {
				// Write as base64
				base64Value := base64.URLEncoding.
					EncodeToString(e.Content[len(BinPrefix):])
				if _, err = buf.Write([]byte("b64:" +
					base64Value)); chk.E(err) {
					return
				}

			} else {
				if err = WriteText(buf, e.Content); chk.E(err) {
					return
				}
			}
		case TAG:
			for t, v := range e.Tags {
				if err = WriteText(buf, v.key); chk.E(err) {
					return
				}
				buf.WriteByte(':')
				// Write the value
				if isBinary(v.value) {
					// Write as base64
					base64Value := base64.URLEncoding.
						EncodeToString(v.value[len(BinPrefix):])
					if _, err = buf.Write([]byte("b64:" +
						base64Value)); chk.E(err) {
						return
					}
				} else {
					// Write plain text
					if err = WriteText(buf, v.value); chk.E(err) {
						return
					}
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
			base64.RawURLEncoding.Encode(b, e.Signature)
			buf.Write(b)
		}
	}
	data = buf.Bytes()
	return
}

func isBinary(data []byte) bool { return bytes.HasPrefix(data, BinPrefix) }

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
