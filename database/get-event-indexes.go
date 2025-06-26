package database

import (
	"bytes"

	"manifold.mleku.dev/chk"
	"manifold.mleku.dev/database/indexes"
	"manifold.mleku.dev/database/indexes/types/fullid"
	"manifold.mleku.dev/database/indexes/types/identhash"
	"manifold.mleku.dev/database/indexes/types/idhash"
	"manifold.mleku.dev/database/indexes/types/number"
	"manifold.mleku.dev/database/indexes/types/pubhash"
	"manifold.mleku.dev/event"
)

func (d *D) GetEventIndexes(ev *event.E) (indices [][]byte, ser *number.Uint40, err error) {

	ser = new(number.Uint40)
	var s uint64
	if s, err = d.Serial(); chk.E(err) {
		return
	}
	if err = ser.Set(s); chk.E(err) {
		panic("there is more than 2^40 events in the database now. " +
			"database needs to be re-consolidated " +
			"(this is ~1,000,000,000,000 at average size 512bytes or 500 terabytes)")
	}

	id := idhash.New()
	var idb []byte
	if idb, err = ev.Id(); chk.E(err) {
		return
	}
	if err = id.FromId(idb); chk.E(err) {
		return
	}
	evIDB := new(bytes.Buffer)
	if err = indexes.IdEnc(id, ser).MarshalWrite(evIDB); chk.E(err) {
		return
	}
	indices = append(indices, evIDB.Bytes())

	fid := fullid.New()
	if err = fid.FromId(idb); chk.E(err) {
		return
	}
	p := pubhash.New()
	if err = p.FromPubkey(ev.Pubkey); chk.E(err) {
		return
	}
	ts := new(number.Uint64)
	ts.Set(uint64(ev.Timestamp))
	evIFiB := new(bytes.Buffer)
	if err = indexes.IdPubkeyTimestampEnc(ser, fid, p, ts).MarshalWrite(evIFiB); chk.E(err) {
		return
	}
	indices = append(indices, evIFiB.Bytes())

	evIPkCaB := new(bytes.Buffer)
	if err = indexes.PubkeyTimestampEnc(p, ts, ser).MarshalWrite(evIPkCaB); chk.E(err) {
		return
	}
	indices = append(indices, evIPkCaB.Bytes())

	evICaB := new(bytes.Buffer)
	if err = indexes.TimestampEnc(ts, ser).MarshalWrite(evICaB); chk.E(err) {
		return
	}
	indices = append(indices, evICaB.Bytes())
	if ev.Tags != nil {
		for _, t := range *ev.Tags {
			k, v := identhash.New(), identhash.New()
			if err = k.FromIdent(t.Key); chk.E(err) {
				return
			}
			if err = v.FromIdent(t.Value); chk.E(err) {
				return
			}
			tb := new(bytes.Buffer)
			if err = indexes.TagTimestampEnc(k, v, ts, ser).MarshalWrite(tb); chk.E(err) {
				return
			}
			indices = append(indices, tb.Bytes())

			ptb := new(bytes.Buffer)
			if err = indexes.PubkeyTagTimestampEnc(p, k, v, ts, ser).MarshalWrite(ptb); chk.E(err) {
				return
			}
			indices = append(indices, ptb.Bytes())
		}
	}

	return
}
