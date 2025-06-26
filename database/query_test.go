package database

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"manifold.mleku.dev/chk"
	"manifold.mleku.dev/event"
	"manifold.mleku.dev/filter"
	"manifold.mleku.dev/p256k"
)

// TestQueryEvents tests the QueryEvents function with various filter combinations.
func TestQueryEvents(t *testing.T) {
	// Create a temporary directory for the database
	tempDir, err := os.MkdirTemp("", "manifold-test-db")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new database
	db := New()
	if err := db.Init(tempDir); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Generate test events
	events, err := generateTestEvents(10)
	if err != nil {
		t.Fatalf("Failed to generate test events: %v", err)
	}

	// Store events in the database
	for _, ev := range events {
		if err := db.StoreEvent(ev); err != nil {
			t.Fatalf("Failed to store event: %v", err)
		}
	}

	// Test cases
	t.Run("FilterByIds", func(t *testing.T) {
		testFilterByIds(t, db, events)
	})

	t.Run("FilterByAuthors", func(t *testing.T) {
		testFilterByAuthors(t, db, events)
	})

	t.Run("FilterByTags", func(t *testing.T) {
		testFilterByTags(t, db, events)
	})

	t.Run("FilterByAuthorsAndTags", func(t *testing.T) {
		testFilterByAuthorsAndTags(t, db, events)
	})

	t.Run("FilterByTimestamp", func(t *testing.T) {
		testFilterByTimestamp(t, db, events)
	})

	t.Run("FilterWithNoSpecificCriteria", func(t *testing.T) {
		testFilterWithNoSpecificCriteria(t, db, events)
	})

	t.Run("SortingAscending", func(t *testing.T) {
		testSortingAscending(t, db, events)
	})

	t.Run("SortingDescending", func(t *testing.T) {
		testSortingDescending(t, db, events)
	})
}

// generateTestEvents generates a set of test events with various properties.
func generateTestEvents(count int) ([]*event.E, error) {
	events := make([]*event.E, count)

	// Create signers for different authors
	signers := make([]*p256k.Signer, 3)
	for i := range signers {
		signers[i] = new(p256k.Signer)
		if err := signers[i].Generate(); chk.E(err) {
			return nil, err
		}
	}

	// Create events with different properties
	for i := 0; i < count; i++ {
		// Use different authors
		authorIndex := i % len(signers)
		signer := signers[authorIndex]

		// Create event
		ev := &event.E{
			Pubkey:    signer.Pub(),
			Timestamp: time.Now().Add(time.Duration(i) * time.Minute).Unix(), // Different timestamps
			Content:   []byte(fmt.Sprintf("Test content %d", i)),
			Tags:      &event.Tags{},
		}

		// Add tags based on index
		if i%2 == 0 {
			*ev.Tags = append(*ev.Tags, event.Tag{
				Key:   []byte("type"),
				Value: []byte("text"),
			})
		}
		if i%3 == 0 {
			*ev.Tags = append(*ev.Tags, event.Tag{
				Key:   []byte("category"),
				Value: []byte("test"),
			})
		}
		if i%4 == 0 {
			*ev.Tags = append(*ev.Tags, event.Tag{
				Key:   []byte("importance"),
				Value: []byte("high"),
			})
		}

		// Sign the event
		if err := ev.Sign(signer); chk.E(err) {
			return nil, err
		}

		events[i] = ev
	}

	return events, nil
}

// testFilterByIds tests filtering events by IDs.
func testFilterByIds(t *testing.T, db *D, events []*event.E) {
	// Get IDs of the first 3 events
	ids := make([][]byte, 3)
	for i := 0; i < 3; i++ {
		var err error
		ids[i], err = events[i].Id()
		if err != nil {
			t.Fatalf("Failed to get event ID: %v", err)
		}
	}

	// Create filter with IDs
	f := filter.F{
		Ids: ids,
	}

	// Query events
	result, err := db.QueryEvents(f)
	if err != nil {
		t.Fatalf("QueryEvents failed: %v", err)
	}

	// Verify results
	if len(result) != len(ids) {
		t.Fatalf("Expected %d events, got %d", len(ids), len(result))
	}

	// Check that all IDs are in the result
	for _, id := range ids {
		found := false
		for _, resultId := range result {
			if bytes.Equal(id, resultId) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("ID %x not found in result", id)
		}
	}
}

// testFilterByAuthors tests filtering events by authors.
func testFilterByAuthors(t *testing.T, db *D, events []*event.E) {
	// Get the pubkey of the first author
	pubkey := events[0].Pubkey

	// Count events with this pubkey
	expectedCount := 0
	for _, ev := range events {
		if bytes.Equal(ev.Pubkey, pubkey) {
			expectedCount++
		}
	}

	// Create filter with author
	f := filter.F{
		Authors: [][]byte{pubkey},
	}

	// Query events
	result, err := db.QueryEvents(f)
	if err != nil {
		t.Fatalf("QueryEvents failed: %v", err)
	}

	// Verify results
	if len(result) != expectedCount {
		t.Fatalf("Expected %d events, got %d", expectedCount, len(result))
	}

	// Verify each result has the correct author
	for _, id := range result {
		ev, err := db.GetEventById(id)
		if err != nil {
			t.Fatalf("Failed to get event: %v", err)
		}
		if !bytes.Equal(ev.Pubkey, pubkey) {
			t.Fatalf("Event has incorrect author")
		}
	}
}

// testFilterByTags tests filtering events by tags.
func testFilterByTags(t *testing.T, db *D, events []*event.E) {
	// Create filter with tag
	f := filter.F{
		Tags: filter.TagMap{
			"type": {[]byte("text")},
		},
	}

	// Count events with this tag
	expectedCount := 0
	for _, ev := range events {
		for _, tag := range *ev.Tags {
			if bytes.Equal(tag.Key, []byte("type")) && bytes.Equal(tag.Value, []byte("text")) {
				expectedCount++
				break
			}
		}
	}

	// Query events
	result, err := db.QueryEvents(f)
	if err != nil {
		t.Fatalf("QueryEvents failed: %v", err)
	}

	// Verify results
	if len(result) != expectedCount {
		t.Fatalf("Expected %d events, got %d", expectedCount, len(result))
	}

	// Verify each result has the correct tag
	for _, id := range result {
		ev, err := db.GetEventById(id)
		if err != nil {
			t.Fatalf("Failed to get event: %v", err)
		}

		found := false
		for _, tag := range *ev.Tags {
			if bytes.Equal(tag.Key, []byte("type")) && bytes.Equal(tag.Value, []byte("text")) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Event does not have the expected tag")
		}
	}
}

// testFilterByAuthorsAndTags tests filtering events by both authors and tags.
func testFilterByAuthorsAndTags(t *testing.T, db *D, events []*event.E) {
	// Get the pubkey of the first author
	pubkey := events[0].Pubkey

	// Create filter with author and tag
	f := filter.F{
		Authors: [][]byte{pubkey},
		Tags: filter.TagMap{
			"type": {[]byte("text")},
		},
	}

	// Count events with this author and tag
	expectedCount := 0
	for _, ev := range events {
		if !bytes.Equal(ev.Pubkey, pubkey) {
			continue
		}

		for _, tag := range *ev.Tags {
			if bytes.Equal(tag.Key, []byte("type")) && bytes.Equal(tag.Value, []byte("text")) {
				expectedCount++
				break
			}
		}
	}

	// Query events
	result, err := db.QueryEvents(f)
	if err != nil {
		t.Fatalf("QueryEvents failed: %v", err)
	}

	// Verify results
	if len(result) != expectedCount {
		t.Fatalf("Expected %d events, got %d", expectedCount, len(result))
	}

	// Verify each result has the correct author and tag
	for _, id := range result {
		ev, err := db.GetEventById(id)
		if err != nil {
			t.Fatalf("Failed to get event: %v", err)
		}

		if !bytes.Equal(ev.Pubkey, pubkey) {
			t.Fatalf("Event has incorrect author")
		}

		found := false
		for _, tag := range *ev.Tags {
			if bytes.Equal(tag.Key, []byte("type")) && bytes.Equal(tag.Value, []byte("text")) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Event does not have the expected tag")
		}
	}
}

// testFilterByTimestamp tests filtering events by timestamp.
func testFilterByTimestamp(t *testing.T, db *D, events []*event.E) {
	// Find the middle timestamp
	timestamps := make([]int64, len(events))
	for i, ev := range events {
		timestamps[i] = ev.Timestamp
	}

	// Sort timestamps
	for i := 0; i < len(timestamps); i++ {
		for j := i + 1; j < len(timestamps); j++ {
			if timestamps[i] > timestamps[j] {
				timestamps[i], timestamps[j] = timestamps[j], timestamps[i]
			}
		}
	}

	// Use the middle timestamp as the since value
	middleIndex := len(timestamps) / 2
	since := timestamps[middleIndex]

	// Create filter with since timestamp
	f := filter.F{
		Since: since,
	}

	// Count events with timestamp >= since
	expectedCount := 0
	for _, ev := range events {
		if ev.Timestamp >= since {
			expectedCount++
		}
	}

	// Query events
	result, err := db.QueryEvents(f)
	if err != nil {
		t.Fatalf("QueryEvents failed: %v", err)
	}

	// Verify results
	if len(result) != expectedCount {
		t.Fatalf("Expected %d events, got %d", expectedCount, len(result))
	}

	// Verify each result has the correct timestamp
	for _, id := range result {
		ev, err := db.GetEventById(id)
		if err != nil {
			t.Fatalf("Failed to get event: %v", err)
		}

		if ev.Timestamp < since {
			t.Fatalf("Event has timestamp %d, which is less than since %d", ev.Timestamp, since)
		}
	}
}

// testFilterWithNoSpecificCriteria tests filtering events with no specific criteria.
func testFilterWithNoSpecificCriteria(t *testing.T, db *D, events []*event.E) {
	// Create empty filter
	f := filter.F{}

	// Query events
	result, err := db.QueryEvents(f)
	if err != nil {
		t.Fatalf("QueryEvents failed: %v", err)
	}

	// Verify results
	if len(result) != len(events) {
		t.Fatalf("Expected %d events, got %d", len(events), len(result))
	}
}

// testSortingAscending tests sorting events in ascending order.
func testSortingAscending(t *testing.T, db *D, events []*event.E) {
	// Create filter with ascending sort
	f := filter.F{
		Sort: "asc",
	}

	// Query events
	result, err := db.QueryEvents(f)
	if err != nil {
		t.Fatalf("QueryEvents failed: %v", err)
	}

	// Verify results
	if len(result) != len(events) {
		t.Fatalf("Expected %d events, got %d", len(events), len(result))
	}

	// Verify sorting
	for i := 1; i < len(result); i++ {
		ev1, err := db.GetEventById(result[i-1])
		if err != nil {
			t.Fatalf("Failed to get event: %v", err)
		}

		ev2, err := db.GetEventById(result[i])
		if err != nil {
			t.Fatalf("Failed to get event: %v", err)
		}

		if ev1.Timestamp > ev2.Timestamp {
			t.Fatalf("Events not sorted in ascending order")
		}
	}
}

// testSortingDescending tests sorting events in descending order.
func testSortingDescending(t *testing.T, db *D, events []*event.E) {
	// Create filter with descending sort
	f := filter.F{
		Sort: "desc",
	}

	// Query events
	result, err := db.QueryEvents(f)
	if err != nil {
		t.Fatalf("QueryEvents failed: %v", err)
	}

	// Verify results
	if len(result) != len(events) {
		t.Fatalf("Expected %d events, got %d", len(events), len(result))
	}

	// Verify sorting
	for i := 1; i < len(result); i++ {
		ev1, err := db.GetEventById(result[i-1])
		if err != nil {
			t.Fatalf("Failed to get event: %v", err)
		}

		ev2, err := db.GetEventById(result[i])
		if err != nil {
			t.Fatalf("Failed to get event: %v", err)
		}

		if ev1.Timestamp < ev2.Timestamp {
			t.Fatalf("Events not sorted in descending order")
		}
	}
}
