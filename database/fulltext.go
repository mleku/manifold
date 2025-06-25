package database

import (
	"bytes"
	"unicode"
	"unicode/utf8"

	"github.com/clipperhouse/uax29/words"

	"manifold.mleku.dev/chk"
	"manifold.mleku.dev/database/indexes"
	"manifold.mleku.dev/database/indexes/types/fulltext"
	"manifold.mleku.dev/database/indexes/types/number"
	"manifold.mleku.dev/event"
	"manifold.mleku.dev/hex"
)

type Words struct {
	ser     *number.Uint40
	ev      *event.E
	wordMap map[string]int
}

// GetFulltextKeys extracts the found words in the Content field of an event,
// and generates indexes that represent the word, zero terminated, the position
// in the sequence of word segments, and appends the provided serial reference
// to the event record.
//
// This process is only performed on events that have a tag "mimetype":"text/*".
func (d *D) GetFulltextKeys(ev *event.E, ser *number.Uint40) (keys [][]byte, err error) {
	mt := ev.Tags.GetFirst([]byte("mimetype"))
	if !bytes.HasPrefix(mt.Value, []byte("text/")) {
		return
	}
	w := d.GetWordsFromContent(ev)
	for i := range w {
		ft := fulltext.New()
		ft.FromWord([]byte(i))
		pos := new(number.Uint24)
		_ = pos.Set(uint32(w[i]))
		buf := new(bytes.Buffer)
		if err = indexes.FullTextWordEnc(ft, pos, ser).MarshalWrite(buf); chk.E(err) {
			return
		}
		keys = append(keys, buf.Bytes())
	}
	return
}

func (d *D) GetWordsFromContent(ev *event.E) (wordMap map[string]int) {
	wordMap = make(map[string]int)
	content := ev.Content
	seg := words.NewSegmenter(content)
	var counter int
	for seg.Next() {
		w := seg.Bytes()
		w = bytes.ToLower(w)
		var ru rune
		ru, _ = utf8.DecodeRune(w)
		// ignore the most common things that aren't words
		if !unicode.IsSpace(ru) &&
			!unicode.IsPunct(ru) &&
			!unicode.IsSymbol(ru) &&
			!bytes.HasSuffix(w, []byte(".jpg")) &&
			!bytes.HasSuffix(w, []byte(".png")) &&
			!bytes.HasSuffix(w, []byte(".jpeg")) &&
			!bytes.HasSuffix(w, []byte(".mp4")) &&
			!bytes.HasSuffix(w, []byte(".mov")) &&
			!bytes.HasSuffix(w, []byte(".aac")) &&
			!bytes.HasSuffix(w, []byte(".mp3")) &&
			!IsEntity(w) &&
			!bytes.Contains(w, []byte(".")) {
			if len(w) == 64 || len(w) == 128 {
				if _, err := hex.Dec(string(w)); err == nil {
					continue
				}
			}
			wordMap[string(w)] = counter
			counter++
		}
	}
	content = content[:0]
	return
}

func IsEntity(w []byte) (is bool) {
	var b []byte
	b = []byte("nostr:")
	if bytes.Contains(w, b) && len(b)+10 < len(w) {
		return true
	}
	b = []byte("npub")
	if bytes.Contains(w, b) && len(b)+5 < len(w) {
		return true
	}
	b = []byte("nsec")
	if bytes.Contains(w, b) && len(b)+5 < len(w) {
		return true
	}
	b = []byte("nevent")
	if bytes.Contains(w, b) && len(b)+5 < len(w) {
		return true
	}
	b = []byte("naddr")
	if bytes.Contains(w, b) && len(b)+5 < len(w) {
		return true
	}
	b = []byte("note")
	if bytes.Contains(w, b) && len(b)+20 < len(w) {
		return true
	}
	b = []byte("lnurl")
	if bytes.Contains(w, b) && len(b)+20 < len(w) {
		return true
	}
	b = []byte("cashu")
	if bytes.Contains(w, b) && len(b)+20 < len(w) {
		return true
	}
	return
}
