package database

import (
	"bytes"
	"time"

	"manifold.mleku.dev/chk"
	"manifold.mleku.dev/database/indexes"
	"manifold.mleku.dev/database/indexes/types/number"
	"manifold.mleku.dev/errorf"
	"manifold.mleku.dev/event"
)

func (d *D) StoreEvent(ev *event.E) (err error) {
	var ev2 *event.E
	var eid []byte
	if eid, err = ev.Id(); chk.E(err) {
		return
	}
	if ev2, err = d.GetEventById(eid); chk.E(err) {
		// so we didn't find it?
	}
	if ev2 != nil {
		// we did found it
		var e2id []byte
		if e2id, err = ev2.Id(); chk.E(err) {
			return
		}
		if bytes.Equal(eid, e2id) {
			err = errorf.E("duplicate event")
			return
		}
	}
	var ser *number.Uint40
	var idxs [][]byte
	if idxs, ser, err = d.GetEventIndexes(ev); chk.E(err) {
		return
	}
	evK := new(bytes.Buffer)
	if err = indexes.EventEnc(ser).MarshalWrite(evK); chk.E(err) {
		return
	}
	ts := new(number.Uint64)
	ts.Set(uint64(time.Now().Unix()))
	// write indexes; none of the above have values.
	for _, v := range idxs {
		if err = d.Set(v, nil); chk.E(err) {
			return
		}
	}
	// event key
	evk := new(bytes.Buffer)
	if err = indexes.EventEnc(ser).MarshalWrite(evk); chk.E(err) {
		return
	}
	// event value (binary encoded)
	evV := new(bytes.Buffer)
	if err = ev.WriteBinary(evV); chk.E(err) {
		return
	}
	if err = d.Set(evk.Bytes(), evV.Bytes()); chk.E(err) {
		return
	}
	return
}
