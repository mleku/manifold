# Manifold Database API Documentation

This document provides an overview of the API functions available in the Manifold database package.

## Table of Contents

1. [Database Initialization](#database-initialization)
2. [Basic Operations](#basic-operations)
3. [Event Operations](#event-operations)
4. [Query Operations](#query-operations)
5. [Logging](#logging)

## Database Initialization

### New

```go
func New() (d *D)
```

Creates a new database instance with default settings.

**Returns:**
- `d *D`: A new database instance

### Init

```go
func (d *D) Init(path string) (err error)
```

Initializes the database with the specified path.

**Parameters:**
- `path string`: The directory path where the database files will be stored

**Returns:**
- `err error`: Any error that occurred during initialization

### Close

```go
func (d *D) Close() (err error)
```

Closes the database.

**Returns:**
- `err error`: Any error that occurred during closing

### Path

```go
func (d *D) Path() string
```

Returns the path where the database files are stored.

**Returns:**
- `string`: The database directory path

## Basic Operations

### Set

```go
func (d *D) Set(k, v []byte) (err error)
```

Sets a key-value pair in the database.

**Parameters:**
- `k []byte`: The key
- `v []byte`: The value (can be nil)

**Returns:**
- `err error`: Any error that occurred

### Get

```go
func (d *D) Get(k []byte) (v []byte, err error)
```

Gets a value by key from the database.

**Parameters:**
- `k []byte`: The key

**Returns:**
- `v []byte`: The value
- `err error`: Any error that occurred

### View

```go
func (d *D) View(fn func(txn *badger.Txn) (err error)) (err error)
```

Executes a read-only transaction.

**Parameters:**
- `fn func(txn *badger.Txn) (err error)`: The function to execute within the transaction

**Returns:**
- `err error`: Any error that occurred

### Update

```go
func (d *D) Update(fn func(txn *badger.Txn) (err error)) (err error)
```

Executes a read-write transaction.

**Parameters:**
- `fn func(txn *badger.Txn) (err error)`: The function to execute within the transaction

**Returns:**
- `err error`: Any error that occurred

### Serial

```go
func (d *D) Serial() (ser uint64, err error)
```

Returns the next monotonic conflict-free unique serial number.

**Returns:**
- `ser uint64`: The next serial number
- `err error`: Any error that occurred

## Event Operations

### StoreEvent

```go
func (d *D) StoreEvent(ev *event.E) (err error)
```

Stores an event in the database, creating all necessary indexes.

**Parameters:**
- `ev *event.E`: The event to store

**Returns:**
- `err error`: Any error that occurred, including "duplicate event" if the event already exists

### GetEventIndexes

```go
func (d *D) GetEventIndexes(ev *event.E) (indices [][]byte, ser *number.Uint40, err error)
```

Generates all the index entries for an event.

**Parameters:**
- `ev *event.E`: The event to generate indexes for

**Returns:**
- `indices [][]byte`: The generated index entries
- `ser *number.Uint40`: The serial number assigned to the event
- `err error`: Any error that occurred

### FindEventSerialById

```go
func (d *D) FindEventSerialById(evId []byte) (ser *number.Uint40, err error)
```

Finds the serial number of an event by its ID.

**Parameters:**
- `evId []byte`: The event ID

**Returns:**
- `ser *number.Uint40`: The serial number of the event
- `err error`: Any error that occurred, including "event not found"

### GetEventFromSerial

```go
func (d *D) GetEventFromSerial(ser *number.Uint40) (ev *event.E, err error)
```

Retrieves an event by its serial number.

**Parameters:**
- `ser *number.Uint40`: The serial number of the event

**Returns:**
- `ev *event.E`: The retrieved event
- `err error`: Any error that occurred

### GetEventById

```go
func (d *D) GetEventById(evId []byte) (ev *event.E, err error)
```

Retrieves an event by its ID.

**Parameters:**
- `evId []byte`: The event ID

**Returns:**
- `ev *event.E`: The retrieved event
- `err error`: Any error that occurred

### GetIdPubkeyTimestampFromSerial

```go
func (d *D) GetIdPubkeyTimestampFromSerial(ser *number.Uint40) (id, pk []byte, ts int64, err error)
```

Retrieves the ID, public key, and timestamp of an event by its serial number.

**Parameters:**
- `ser *number.Uint40`: The serial number of the event

**Returns:**
- `id []byte`: The event ID
- `pk []byte`: The public key
- `ts int64`: The timestamp
- `err error`: Any error that occurred

## Query Operations

### QueryEvents

```go
func (d *D) QueryEvents(f filter.F) (eventIds [][]byte, err error)
```

Queries events based on the provided filter criteria.

**Parameters:**
- `f filter.F`: The filter criteria, which can include:
  - `Ids [][]byte`: Specific event IDs to retrieve *(if this field is present, any others will be invalid and return an error)*
  - `Authors [][]byte`: Author public keys to filter by
  - `Tags map[string][][]byte`: Tags to filter by
  - `Since int64`: Minimum timestamp (inclusive)
  - `Until int64`: Maximum timestamp (inclusive)
  - `NotIds [][]byte`: Event IDs to exclude
  - `NotAuthors [][]byte`: Author public keys to exclude
  - `NotTags map[string][][]byte`: Tags to exclude
  - `Sort string`: Sort order ("asc" or "desc")

**Returns:**
- `eventIds [][]byte`: The IDs of events matching the filter criteria
- `err error`: Any error that occurred

## Logging

### NewLogger

```go
func NewLogger(logLevel int, label string) (l *logger)
```

Creates a new logger for the database.

**Parameters:**
- `logLevel int`: The initial log level
- `label string`: A label for the logger

**Returns:**
- `l *logger`: The new logger

### SetLogLevel

```go
func (l *logger) SetLogLevel(level int)
```

Sets the log level of the logger.

**Parameters:**
- `level int`: The new log level

### Log Methods

The logger provides the following methods for logging at different levels:

```go
func (l *logger) Errorf(s string, i ...interface{})
func (l *logger) Warningf(s string, i ...interface{})
func (l *logger) Infof(s string, i ...interface{})
func (l *logger) Debugf(s string, i ...interface{})
```

Each method takes a format string and optional arguments, similar to `fmt.Printf`.