// Package text is a reader and writer that processes received raw text, and
// escapes it using the simple rule that messages are separated by line breaks,
// and so only linebreaks in the text need to be escaped. Because of the need of
// the escape backslash, there is a second escaping rule that means all single
// backslashes must be turned into double, and for both newline and backslash,
// they are returned to the raw form as a single backslash and a newline.
//
// This was designed this way after dealing with the multiplicity of escape
// conventions used by JSON, which are incompatible with the requirement that
// only one marshaled form is valid, because it must hash to generate the Id,
// which is then signed on to authenticate the message.
//
// With this convention, there is no ambiguity or ways to misinterpret the data,
// simplifying the processing of text fields in the data structure (content and
// tags).
//
// A side note: where a field represents purely binary data, it must be prefixed
// by `b64:` and as such this text in the first 4 characters of a string is a
// reserved word and cannot appear in any valid *text* data.
package text

import (
	"bytes"
	"io"

	"manifold.mleku.dev/chk"
)

func Write(w io.Writer, b []byte) (err error) {
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

func Read(r io.Reader) (b []byte, err error) {
	buf := new(bytes.Buffer)
	var inEscape bool
	var n int
	rb := make([]byte, 1)
	for {
		if n, err = r.Read(rb); chk.E(err) {
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
