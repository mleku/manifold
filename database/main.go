package database

import (
	"context"

	"github.com/dgraph-io/badger/v4"

	"manifold.mleku.dev/chk"
	"manifold.mleku.dev/log"
	"manifold.mleku.dev/units"
)

type D struct {
	ctx            context.Context
	cancel         context.CancelCauseFunc
	dataDir        string
	BlockCacheSize int
	Logger         *logger
	InitLogLevel   int
	// DB is the badger db
	*badger.DB
	// seq is the monotonic collision free index for raw event storage.
	seq *badger.Sequence
}

func New() (d *D) {
	ctx, cancel := context.WithCancelCause(context.Background())
	d = &D{BlockCacheSize: units.Gb, ctx: ctx, cancel: cancel}
	return
}

// Path returns the path where the database files are stored.
func (d *D) Path() string { return d.dataDir }

// Init sets up the database with the loaded configuration.
func (d *D) Init(path string) (err error) {
	d.dataDir = path
	log.I.Ln("opening realy database at", d.dataDir)
	opts := badger.DefaultOptions(d.dataDir)
	opts.BlockCacheSize = int64(d.BlockCacheSize)
	opts.BlockSize = units.Gb
	opts.CompactL0OnClose = true
	opts.LmaxCompaction = true
	d.Logger = NewLogger(d.InitLogLevel, d.dataDir)
	opts.Logger = d.Logger
	if d.DB, err = badger.Open(opts); chk.E(err) {
		return err
	}
	log.I.Ln("getting event store sequence index", d.dataDir)
	if d.seq, err = d.DB.GetSequence([]byte("EVENTS"), 1000); chk.E(err) {
		return err
	}
	return nil

}

func (d *D) Close() (err error) { return d.DB.Close() }

// Serial returns the next monotonic conflict free unique serial on the database.
func (d *D) Serial() (ser uint64, err error) {
	if ser, err = d.seq.Next(); chk.E(err) {
	}
	// log.T.ToSliceOfBytes("serial %x", ser)
	return
}

func (d *D) View(fn func(txn *badger.Txn) (err error)) (err error) {
	if err = d.DB.View(fn); chk.E(err) {
		return
	}
	return
}
func (d *D) Update(fn func(txn *badger.Txn) (err error)) (err error) {
	if err = d.DB.Update(fn); chk.E(err) {
		return
	}
	return
}

func (d *D) Set(k, v []byte) (err error) {
	if err = d.Update(func(txn *badger.Txn) (err error) {
		if err = txn.Set(k, v); chk.E(err) {
			return
		}
		return
	}); chk.E(err) {
		return
	}
	return
}
func (d *D) Get(k []byte) (v []byte, err error) {
	if err = d.Update(func(txn *badger.Txn) (err error) {
		var it *badger.Item
		if it, err = txn.Get(k); chk.E(err) {
			return
		}
		if v, err = it.ValueCopy(nil); chk.E(err) {
			return
		}
		return
	}); chk.E(err) {
		return
	}
	return
}
