# Manifold Event Package Documentation

This document provides an overview of the event package in Manifold, which implements a structured event system for data exchange and storage.

## Table of Contents

1. [Overview](#overview)
2. [Event Structure](#event-structure)
3. [Tags](#tags)
4. [Serialization](#serialization)
5. [Event Operations](#event-operations)

## Overview

The event package provides a structured way to create, serialize, sign, and verify events in the Manifold system. Events are the primary data structure used for exchanging information, with support for both text and binary content, along with flexible metadata through tags.

## Event Structure

### E (Event)

The `E` struct represents a Manifold event with the following fields:
- `Pubkey`: The public key of the event creator
- `Timestamp`: The creation time of the event
- `Content`: The main content of the event (can be text or binary)
- `Tags`: Optional metadata tags associated with the event
- `Signature`: Cryptographic signature to verify the event's authenticity

```
type E struct {
    Pubkey    []byte
    Timestamp int64
    Content   []byte
    Tags      *Tags
    Signature []byte
}
```

Unlike some other event systems, Manifold events do not include a "kind" field, as this functionality is handled through tags. The event ID is derived from the event content (excluding the signature) rather than being stored explicitly.

## Tags

### Tag

Tags provide a flexible way to annotate events with metadata. Each tag consists of a key-value pair.

```
type Tag struct {
    Key   []byte
    Value []byte
}
```

### Tags

A collection of Tag objects associated with an event.

```
type Tags []Tag
```

### Tag Methods

#### GetAll

Returns all tags with the specified key.

```
func (t Tags) GetAll(k []byte) (tt Tags)
```

**Parameters:**
- `k []byte`: The key to search for

**Returns:**
- `tt Tags`: All tags matching the key

#### GetFirst

Returns the first tag with the specified key. Useful when a tag is expected to appear only once.

```
func (t Tags) GetFirst(k []byte) (tt Tag)
```

**Parameters:**
- `k []byte`: The key to search for

**Returns:**
- `tt Tag`: The first tag matching the key

## Serialization

Events are serialized in a specific format with sentinel markers for each field:

```
PUBKEY:<base64-encoded-pubkey>
TIMESTAMP:<encoded-timestamp>
CONTENT:<text-or-base64-encoded-binary>
TAG:<key>:<value>
TAG:<key>:<value>
...
SIGNATURE:<base64-encoded-signature>
```

Binary content is prefixed with "BIN:" and encoded as base64 with a "b64:" prefix in the serialized form.

### Marshal

Serializes an event to its byte representation.

```
func (e *E) Marshal() (data []byte, err error)
```

**Returns:**
- `data []byte`: The serialized event
- `err error`: Any error that occurred during marshaling

### Unmarshal

Deserializes a byte representation into an event.

```
func (e *E) Unmarshal(data []byte) (err error)
```

**Parameters:**
- `data []byte`: The serialized event data

**Returns:**
- `err error`: Any error that occurred during unmarshaling

## Event Operations

### Id

Calculates the ID of an event, which is the SHA-256 hash of the event's canonical representation (without signature).

```
func (e *E) Id() (id []byte, err error)
```

**Returns:**
- `id []byte`: The event ID
- `err error`: Any error that occurred

### Sign

Signs an event using the provided signer.

```
func (e *E) Sign(sign signer.I) (err error)
```

**Parameters:**
- `sign signer.I`: The signer interface implementation

**Returns:**
- `err error`: Any error that occurred during signing

### Verify

Verifies the signature of an event.

```
func (e *E) Verify() (valid bool, err error)
```

**Returns:**
- `valid bool`: Whether the signature is valid
- `err error`: Any error that occurred during verification
