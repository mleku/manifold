package database

import (
	"bytes"
	"github.com/dgraph-io/badger/v4"
	"sort"

	"manifold.mleku.dev/chk"
	"manifold.mleku.dev/database/indexes"
	"manifold.mleku.dev/database/indexes/types/identhash"
	"manifold.mleku.dev/database/indexes/types/number"
	"manifold.mleku.dev/database/indexes/types/pubhash"
	"manifold.mleku.dev/filter"
)

// QueryEvents finds events that match the given filter and returns their IDs.
// The results are sorted according to the Sort field in the filter.
func (d *D) QueryEvents(f filter.F) (eventIds [][]byte, err error) {
	// If specific IDs are provided, just return those
	if len(f.Ids) > 0 {
		return f.Ids, nil
	}

	// Create a map to store unique event serials
	eventSerials := make(map[uint64]struct{})

	// Use View transaction to read from the database
	if err = d.View(func(txn *badger.Txn) (err error) {
		// If both authors and tags are specified, use the PubkeyTagTimestamp index
		if len(f.Authors) > 0 && len(f.Tags) > 0 {
			for _, author := range f.Authors {
				p := pubhash.New()
				if err = p.FromPubkey(author); chk.E(err) {
					return
				}

				for tagName, tagValues := range f.Tags {
					for _, tagValue := range tagValues {
						k, v := identhash.New(), identhash.New()
						if err = k.FromIdent([]byte(tagName)); chk.E(err) {
							return
						}
						if err = v.FromIdent(tagValue); chk.E(err) {
							return
						}

						// Create timestamp range
						tsStart := new(number.Uint64)
						tsEnd := new(number.Uint64)

						if f.Since > 0 {
							tsStart.Set(uint64(f.Since))
						}
						if f.Until > 0 {
							tsEnd.Set(uint64(f.Until))
						} else {
							tsEnd.Set(^uint64(0)) // Max value if Until not specified
						}

						// Create prefix for PubkeyTagTimestamp index
						prefix := new(bytes.Buffer)
						if err = indexes.PubkeyTagTimestampEnc(p, k, v, nil, nil).MarshalWrite(prefix); chk.E(err) {
							return
						}

						// Iterate over events with this author and tag
						it := txn.NewIterator(badger.IteratorOptions{Prefix: prefix.Bytes()})
						defer it.Close()

						for it.Seek(prefix.Bytes()); it.Valid(); it.Next() {
							item := it.Item()
							k := item.KeyCopy(nil)
							buf := bytes.NewBuffer(k)

							// Decode the key
							p2, k2, v2, ts, ser := indexes.PubkeyTagTimestampVars()
							if err = indexes.PubkeyTagTimestampDec(p2, k2, v2, ts, ser).UnmarshalRead(buf); chk.E(err) {
								return
							}

							// Check timestamp range
							tsValue := ts.Get()
							if (f.Since <= 0 || int64(tsValue) >= f.Since) &&
								(f.Until <= 0 || int64(tsValue) <= f.Until) {
								// Add to results
								eventSerials[ser.Get()] = struct{}{}
							}
						}
					}
				}
			}
			// If only authors are specified, use the PubkeyTimestamp index
		} else if len(f.Authors) > 0 {
			for _, author := range f.Authors {
				p := pubhash.New()
				if err = p.FromPubkey(author); chk.E(err) {
					return
				}

				// Create timestamp range
				tsStart := new(number.Uint64)
				tsEnd := new(number.Uint64)

				if f.Since > 0 {
					tsStart.Set(uint64(f.Since))
				}
				if f.Until > 0 {
					tsEnd.Set(uint64(f.Until))
				} else {
					tsEnd.Set(^uint64(0)) // Max value if Until not specified
				}

				// Create prefix for PubkeyTimestamp index
				prefix := new(bytes.Buffer)
				if err = indexes.PubkeyTimestampEnc(p, nil, nil).MarshalWrite(prefix); chk.E(err) {
					return
				}

				// Iterate over events with this author
				it := txn.NewIterator(badger.IteratorOptions{Prefix: prefix.Bytes()})
				defer it.Close()

				for it.Seek(prefix.Bytes()); it.Valid(); it.Next() {
					item := it.Item()
					k := item.KeyCopy(nil)
					buf := bytes.NewBuffer(k)

					// Decode the key
					p2, ts, ser := indexes.PubkeyTimestampVars()
					if err = indexes.PubkeyTimestampDec(p2, ts, ser).UnmarshalRead(buf); chk.E(err) {
						return
					}

					// Check timestamp range
					tsValue := ts.Get()
					if (f.Since <= 0 || int64(tsValue) >= f.Since) &&
						(f.Until <= 0 || int64(tsValue) <= f.Until) {
						// Add to results
						eventSerials[ser.Get()] = struct{}{}
					}
				}
			}
		} else if len(f.Tags) > 0 {
			// If tags are specified, use the TagTimestamp index
			for tagName, tagValues := range f.Tags {
				for _, tagValue := range tagValues {
					k, v := identhash.New(), identhash.New()
					if err = k.FromIdent([]byte(tagName)); chk.E(err) {
						return
					}
					if err = v.FromIdent(tagValue); chk.E(err) {
						return
					}

					// Create timestamp range
					tsStart := new(number.Uint64)
					tsEnd := new(number.Uint64)

					if f.Since > 0 {
						tsStart.Set(uint64(f.Since))
					}
					if f.Until > 0 {
						tsEnd.Set(uint64(f.Until))
					} else {
						tsEnd.Set(^uint64(0)) // Max value if Until not specified
					}

					// Create prefix for TagTimestamp index
					prefix := new(bytes.Buffer)
					if err = indexes.TagTimestampEnc(k, v, nil, nil).MarshalWrite(prefix); chk.E(err) {
						return
					}

					// Iterate over events with this tag
					it := txn.NewIterator(badger.IteratorOptions{Prefix: prefix.Bytes()})
					defer it.Close()

					for it.Seek(prefix.Bytes()); it.Valid(); it.Next() {
						item := it.Item()
						k := item.KeyCopy(nil)
						buf := bytes.NewBuffer(k)

						// Decode the key
						k2, v2, ts, ser := indexes.TagTimestampVars()
						if err = indexes.TagTimestampDec(k2, v2, ts, ser).UnmarshalRead(buf); chk.E(err) {
							return
						}

						// Check timestamp range
						tsValue := ts.Get()
						if (f.Since <= 0 || int64(tsValue) >= f.Since) &&
							(f.Until <= 0 || int64(tsValue) <= f.Until) {
							// Add to results
							eventSerials[ser.Get()] = struct{}{}
						}
					}
				}
			}
		} else if f.Since > 0 || f.Until > 0 {
			// If only timestamp range is specified, use the Timestamp index
			tsStart := new(number.Uint64)
			tsEnd := new(number.Uint64)

			if f.Since > 0 {
				tsStart.Set(uint64(f.Since))
			}
			if f.Until > 0 {
				tsEnd.Set(uint64(f.Until))
			} else {
				tsEnd.Set(^uint64(0)) // Max value if Until not specified
			}

			// Create prefix for Timestamp index
			prefix := new(bytes.Buffer)
			if err = indexes.TimestampEnc(nil, nil).MarshalWrite(prefix); chk.E(err) {
				return
			}

			// Iterate over events in the timestamp range
			it := txn.NewIterator(badger.IteratorOptions{Prefix: prefix.Bytes()})
			defer it.Close()

			for it.Seek(prefix.Bytes()); it.Valid(); it.Next() {
				item := it.Item()
				k := item.KeyCopy(nil)
				buf := bytes.NewBuffer(k)

				// Decode the key
				ts, ser := indexes.TimestampVars()
				if err = indexes.TimestampDec(ts, ser).UnmarshalRead(buf); chk.E(err) {
					return
				}

				// Check timestamp range
				tsValue := ts.Get()
				if (f.Since <= 0 || int64(tsValue) >= f.Since) &&
					(f.Until <= 0 || int64(tsValue) <= f.Until) {
					// Add to results
					eventSerials[ser.Get()] = struct{}{}
				}
			}
		} else {
			// If no specific criteria, use the Event index to get all events
			prefix := new(bytes.Buffer)
			if err = indexes.EventEnc(nil).MarshalWrite(prefix); chk.E(err) {
				return
			}

			// Iterate over all events
			it := txn.NewIterator(badger.IteratorOptions{Prefix: prefix.Bytes()})
			defer it.Close()

			for it.Seek(prefix.Bytes()); it.Valid(); it.Next() {
				item := it.Item()
				k := item.KeyCopy(nil)
				buf := bytes.NewBuffer(k)

				// Decode the key
				ser := indexes.EventVars()
				if err = indexes.EventDec(ser).UnmarshalRead(buf); chk.E(err) {
					return
				}

				// Add to results
				eventSerials[ser.Get()] = struct{}{}
			}
		}

		return nil
	}); chk.E(err) {
		return nil, err
	}

	// Convert serials to event IDs
	serials := make([]uint64, 0, len(eventSerials))
	for serial := range eventSerials {
		serials = append(serials, serial)
	}

	// Sort serials based on the Sort field
	if f.Sort == "desc" {
		sort.Slice(serials, func(i, j int) bool {
			return serials[i] > serials[j]
		})
	} else {
		// Default to ascending order
		sort.Slice(serials, func(i, j int) bool {
			return serials[i] < serials[j]
		})
	}

	// Get event IDs from serials
	for _, serial := range serials {
		ser := new(number.Uint40)
		if err = ser.Set(serial); chk.E(err) {
			return nil, err
		}

		var id []byte
		if id, err = d.GetEventIdFromSerial(ser); chk.E(err) {
			return nil, err
		}

		eventIds = append(eventIds, id)
	}

	return eventIds, nil
}
