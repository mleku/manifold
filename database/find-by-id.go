package database

import (
	"bytes"
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"manifold.mleku.dev/chk"
	"manifold.mleku.dev/database/indexes"
	"manifold.mleku.dev/database/indexes/types/idhash"
	"manifold.mleku.dev/database/indexes/types/number"
	"manifold.mleku.dev/event"
)

func (d *D) FindEventSerialById(evId []byte) (ser *number.Uint40, err error) {
	id := idhash.New()
	if err = id.FromId(evId); chk.E(err) {
		return
	}
	// find by id
	if err = d.View(func(txn *badger.Txn) (err error) {
		key := new(bytes.Buffer)
		if err = indexes.IdSearch(id).MarshalWrite(key); chk.E(err) {
			return
		}
		it := txn.NewIterator(badger.IteratorOptions{Prefix: key.Bytes()})
		defer it.Close()
		for it.Seek(key.Bytes()); it.Valid(); it.Next() {
			item := it.Item()
			k := item.KeyCopy(nil)
			buf := bytes.NewBuffer(k)
			ser = new(number.Uint40)
			if err = indexes.IdDec(id, ser).UnmarshalRead(buf); chk.E(err) {
				return
			}
		}
		return
	}); err != nil {
		return
	}
	if ser == nil {
		err = fmt.Errorf("event %0x not found", evId)
		return
	}
	return
}

func (d *D) GetEventFromSerial(ser *number.Uint40) (ev *event.E, err error) {
	if err = d.View(func(txn *badger.Txn) (err error) {
		enc := indexes.EventEnc(ser)
		kb := new(bytes.Buffer)
		if err = enc.MarshalWrite(kb); chk.E(err) {
			return
		}
		var item *badger.Item
		if item, err = txn.Get(kb.Bytes()); err != nil {
			return
		}
		var val []byte
		if val, err = item.ValueCopy(nil); chk.E(err) {
			return
		}
		var e event.E
		vr := bytes.NewBuffer(val)
		if err = ev.ReadBinary(vr); chk.E(err) {
			return
		}
		ev = &e
		return
	}); err != nil {
		return
	}
	return
}

func (d *D) GetEventIdFromSerial(ser *number.Uint40) (id []byte, err error) {
	if err = d.View(func(txn *badger.Txn) (err error) {
		enc := indexes.IdPubkeyTimestampSearch(ser)
		prf := new(bytes.Buffer)
		if err = enc.MarshalWrite(prf); chk.E(err) {
			return
		}
		it := txn.NewIterator(badger.IteratorOptions{Prefix: prf.Bytes()})
		defer it.Close()
		for it.Seek(prf.Bytes()); it.Valid(); it.Next() {
			item := it.Item()
			key := item.KeyCopy(nil)
			kbuf := bytes.NewBuffer(key)
			_, t, p, ca := indexes.IdPubkeyTimestampVars()
			dec := indexes.IdPubkeyTimestampDec(ser, t, p, ca)
			if err = dec.UnmarshalRead(kbuf); chk.E(err) {
				return
			}
			id = t.Bytes()
		}
		return
	}); chk.E(err) {
		return
	}
	return
}

func (d *D) GetEventById(evId []byte) (ev *event.E, err error) {
	var ser *number.Uint40
	if ser, err = d.FindEventSerialById(evId); chk.E(err) {
		return
	}
	if ev, err = d.GetEventFromSerial(ser); chk.E(err) {
		return
	}
	return
}
